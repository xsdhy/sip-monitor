package rtcp

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sip-monitor/src/pkg/hep"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type RTCPReportService struct {
	mu         sync.RWMutex
	RTCPReport map[string]*CallRTCPReports
	cleanTimer *time.Ticker  // 定时器
	stopCh     chan struct{} // 停止信号
	logger     *logrus.Logger
}

func NewRTCPReportService(logger *logrus.Logger) *RTCPReportService {
	service := &RTCPReportService{
		RTCPReport: make(map[string]*CallRTCPReports),
		stopCh:     make(chan struct{}),
		logger:     logger,
	}

	// 启动定时清理任务，每分钟执行一次
	service.cleanTimer = time.NewTicker(1 * time.Minute)
	go service.cleanupRoutine()

	return service
}

// 停止服务
func (s *RTCPReportService) Stop() {
	if s.cleanTimer != nil {
		s.cleanTimer.Stop()
		close(s.stopCh)
	}
}

// 定时清理任务
func (s *RTCPReportService) cleanupRoutine() {
	for {
		select {
		case <-s.cleanTimer.C:
			s.cleanupOldReports()
		case <-s.stopCh:
			return
		}
	}
}

// 清理旧报告，使用分段处理减少锁定时间
func (s *RTCPReportService) cleanupOldReports() {
	cutoffTime := time.Now().Add(-3 * time.Minute)

	// 第一步：获取需要删除的callID列表，减少锁定时间
	var expiredCallIDs []string

	s.mu.RLock()
	for callID, report := range s.RTCPReport {
		if report.LastUpdated.Before(cutoffTime) {
			expiredCallIDs = append(expiredCallIDs, callID)
		}
	}
	s.mu.RUnlock()

	// 如果没有过期数据，直接返回
	if len(expiredCallIDs) == 0 {
		return
	}

	// 第二步：批量删除过期记录，每批次最多删除100个
	const batchSize = 100
	for i := 0; i < len(expiredCallIDs); i += batchSize {
		end := i + batchSize
		if end > len(expiredCallIDs) {
			end = len(expiredCallIDs)
		}

		currentBatch := expiredCallIDs[i:end]

		// 获取写锁并删除当前批次
		s.mu.Lock()
		for _, callID := range currentBatch {
			// 再次检查是否仍然过期（可能在获取锁期间被更新）
			report, exists := s.RTCPReport[callID]
			if exists && report.LastUpdated.Before(cutoffTime) {
				delete(s.RTCPReport, callID)
			}
		}
		s.mu.Unlock()

		// 释放CPU，让其他goroutine有机会执行
		if end < len(expiredCallIDs) {
			runtime.Gosched()
		}
	}
}

func (s *RTCPReportService) GetCallRTCPReportByCallID(callID string) *CallRTCPReports {
	s.mu.RLock()
	defer s.mu.RUnlock()

	report, ok := s.RTCPReport[callID]
	if !ok {
		return nil
	}
	return report
}

// AddLegRTCPReport 添加或更新一个Leg的RTCP报告
func (s *RTCPReportService) AddLegRTCPReport(ip string, callID string, hepMsg *hep.HepMsg, rawReport *RTCPPacket) {
	s.mu.Lock()
	defer s.mu.Unlock()

	report, ok := s.RTCPReport[callID]
	if !ok {
		report = &CallRTCPReports{
			CallID:      callID,
			Legs:        make(map[string]*LegRTCPReport),
			LastUpdated: time.Now(),
		}
		s.RTCPReport[callID] = report
	} else {
		// 更新最后更新时间
		report.LastUpdated = time.Now()
	}

	legRepot, ok := report.Legs[hepMsg.IP4SourceAddress]
	if !ok {
		legRepot = &LegRTCPReport{
			NodeIP:     ip,
			SrcAddr:    hepMsg.IP4SourceAddress,
			SrcPort:    hepMsg.SourcePort,
			DstAddr:    hepMsg.IP4DestinationAddress,
			DstPort:    hepMsg.DestinationPort,
			RawPackets: make([]*RTCPPacket, 0),
		}
		report.Legs[hepMsg.IP4SourceAddress] = legRepot
		legRepot.RawPackets = append(legRepot.RawPackets, rawReport)
		return
	}
	legRepot.RawPackets = append(legRepot.RawPackets, rawReport)
}

// 清理指定callID的报告
func (s *RTCPReportService) ClearCallReport(callID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.RTCPReport, callID)
}

// 接收RTCP包（HepServer 接受到的包会异步调用这个方法）
func (s *RTCPReportService) ReceiveRTCPPacket(ip string, hepMsg *hep.HepMsg) error {
	if hepMsg.ProtocolType != hep.ProtocolTypeRTCP {
		return nil
	}
	if hepMsg.InternalCorrelationID == "" {
		return nil
	}
	direction := fmt.Sprintf("%s:%d-%s:%d", hepMsg.IP4SourceAddress, hepMsg.SourcePort, hepMsg.IP4DestinationAddress, hepMsg.DestinationPort)

	var rtcpPacket RTCPPacket
	err := json.Unmarshal(hepMsg.Body, &rtcpPacket)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"callID":    hepMsg.InternalCorrelationID,
			"direction": direction,
			"packet":    string(hepMsg.Body),
			"error":     err,
		}).Error("unmarshal rtcp packet failed")
		return err
	}

	rtcpPacket.Raw = string(hepMsg.Body)
	rtcpPacket.TimestampMicro = time.Unix(int64(hepMsg.Timestamp), 0).Add(time.Microsecond * time.Duration(hepMsg.TimestampMicro)).UnixMicro()

	s.AddLegRTCPReport(ip, hepMsg.InternalCorrelationID, hepMsg, &rtcpPacket)
	return nil
}
