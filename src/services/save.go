package services

import (
	"context"
	"sync"
	"time"

	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"sip-monitor/src/pkg/rtcp"
	"sip-monitor/src/pkg/util"

	"github.com/sirupsen/logrus"
)

type SaveService struct {
	logger          *logrus.Logger
	repository      model.Repository
	callRecordCache map[string]*entity.Call
	cacheMutex      sync.RWMutex
	SaveToDBQueue   chan entity.SIP
	rtcpService     *rtcp.RTCPReportService
}

func NewSaveService(logger *logrus.Logger, repository model.Repository, rtcpService *rtcp.RTCPReportService) *SaveService {
	s := &SaveService{
		logger:          logger,
		repository:      repository,
		callRecordCache: make(map[string]*entity.Call),
		cacheMutex:      sync.RWMutex{},
		SaveToDBQueue:   make(chan entity.SIP, 20000),
		rtcpService:     rtcpService,
	}
	s.InitSaveToDBRunner()
	// 启动处理队列的任务
	go s.SaveToDBRunner()
	return s
}

// 定时将缓存刷新到数据库
func (s *SaveService) InitSaveToDBRunner() {
	// 启动周期性刷新缓存到数据库的任务
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			s.FlushCacheToDB()
		}
	}()
}

// 将缓存刷新到数据库
func (s *SaveService) FlushCacheToDB() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	now := time.Now()
	count := 0

	// 遍历缓存中的记录
	for callID, record := range s.callRecordCache {
		timeoutDuration := 1 * time.Minute
		if record.CallStatus == 2 {
			// 通话已建立，15分钟未结束，则认为通话已结束，其他情况为1分钟超时
			timeoutDuration = 15 * time.Minute
		}

		// 检查记录是否已完成或超过一定时间未更新
		if record.EndTime != nil ||
			(record.CreateTime != nil && now.Sub(*record.CreateTime) > timeoutDuration) {
			// 将记录保存到数据库
			ctx := context.Background()

			// 计算通话持续时间
			if record.CreateTime != nil {
				if record.EndTime != nil {
					record.CallDuration = int(record.EndTime.Sub(*record.CreateTime) / time.Second)
				}
				if record.RingingTime != nil {
					record.RingingDuration = int(record.RingingTime.Sub(*record.CreateTime) / time.Second)
				}
				if record.AnswerTime != nil && record.EndTime != nil {
					record.TalkDuration = int(record.EndTime.Sub(*record.AnswerTime) / time.Second)
				}
			}

			go s.dealRTCPReport(callID)

			// 使用内部函数进行更新，便于测试
			err := s.repository.CreateCall(ctx, record)
			if err != nil {
				logrus.WithError(err).Error("更新SIP呼叫记录失败")
			} else {
				count++
				// 从缓存中删除已保存的记录
				delete(s.callRecordCache, callID)
			}
		}
	}

	if count > 0 {
		logrus.WithField("count", count).Info("成功将缓存中的SIP呼叫记录写入数据库")
	}
}

func (s *SaveService) SaveToDBRunner() {
	for item := range s.SaveToDBQueue {
		s.SaveOptimized(item)
	}
}

func (s *SaveService) SaveOptimized(item entity.SIP) {
	// 忽略注册和通知消息
	if item.CSeqMethod == "REGISTER" || item.CSeqMethod == "NOTIFY" {
		return
	}

	// 始终需要在Record表中，新增一条记录
	go func() {
		if item.CallID == "" {
			return
		}
		ctx := context.TODO()
		// 将SIP转换为Record
		record := entity.Record{
			NodeIP:         item.NodeIP,
			SIPCallID:      item.CallID,
			Method:         item.Title,
			ResponseDesc:   item.ResponseDesc,
			ToUser:         item.ToUser,
			FromUser:       item.FromUser,
			SrcAddr:        item.SrcAddr,
			DstAddr:        item.DstAddr,
			CreateTime:     item.CreateTime,
			TimestampMicro: item.TimestampMicro,
		}

		err := s.repository.CreateRecord(ctx, &record)
		if err != nil {
			logrus.WithError(err).Error("保存SIP消息记录失败")
			return
		}

		// 清理Raw文本中的不支持字符
		sanitizedRaw := ""
		if item.Raw != nil {
			sanitizedRaw = util.SanitizeRawText(*item.Raw)
		}

		err = s.repository.CreateRecordRaw(ctx, &entity.RecordRaw{
			ID:         record.ID,
			Raw:        sanitizedRaw,
			CreateTime: item.CreateTime,
		})
		if err != nil {
			logrus.WithError(err).Error("保存SIP消息记录失败")
			return
		}
	}()

	// 使用内存缓存处理呼叫记录
	s.updateCallRecordInCache(item)
}

// 在内存缓存中更新呼叫记录
func (s *SaveService) updateCallRecordInCache(item entity.SIP) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	callID := item.CallID
	record, exists := s.callRecordCache[callID]

	// 如果记录不存在，创建一个新记录
	if !exists {
		// 对于新记录，只有INVITE方法才会创建
		if item.Title == "INVITE" && item.CSeqMethod == "INVITE" {
			record = &entity.Call{
				NodeIP:    item.NodeIP,
				SIPCallID: item.CallID,
				SessionID: item.SessionID,
			}

			record.ToUser = item.ToUser
			record.FromUser = item.FromUser
			record.UserAgent = item.UserAgent
			record.SrcAddr = item.SrcAddr
			record.DstAddr = item.DstAddr
			record.TimestampMicro = item.TimestampMicro
			record.CreateTime = &item.CreateTime

			s.callRecordCache[callID] = record
		}
		return
	}

	// 对于已存在的记录，更新相关字段
	switch item.CSeqMethod {
	case "INVITE", "BYE", "ACK", "CANCEL", "UPDATE":
		switch item.Title {
		case "INVITE":
			// 已经在创建记录时处理
			break
		case "180", "183": // Ringing
			if record.RingingTime == nil {
				ringingTime := item.CreateTime
				record.RingingTime = &ringingTime
				record.CallStatus = 1
			}
		case "200": // OK
			if (item.CSeqMethod == "ACK" || item.CSeqMethod == "INVITE") && record.AnswerTime == nil {
				answerTime := item.CreateTime
				record.AnswerTime = &answerTime
				record.CallStatus = 2
			} else if item.CSeqMethod == "BYE" && record.EndTime == nil {
				endTime := item.CreateTime
				record.EndTime = &endTime
				record.CallStatus = 3
				record.HangupCode = 200
				record.HangupCause = "Normal Clearing"
			}
		case "CANCEL", "480", "487", "404", "403", "408", "413", "416", "486", "488", "513", "500", "503", "504", "580": // Error or Cancel
			if record.EndTime == nil {
				endTime := item.CreateTime
				record.EndTime = &endTime
				record.CallStatus = 3
			}

			// 只有响应的时候，才设置挂断原因
			if item.SrcAddr == record.DstAddr {
				if record.HangupCode == 0 {
					record.HangupCode = item.ResponseCode
				}
				if record.HangupCause == "" {
					record.HangupCause = item.ResponseDesc
				}
			}
		}
	}

	// 如果记录已完成（已结束），立即写入数据库并从缓存移除
	if record.EndTime != nil {
		ctx := context.Background()

		// 计算通话持续时间
		if record.CreateTime != nil {
			record.CallDuration = int(record.EndTime.Sub(*record.CreateTime) / time.Second)
			if record.RingingTime != nil {
				if record.AnswerTime != nil {
					record.RingingDuration = int(record.AnswerTime.Sub(*record.RingingTime) / time.Second)
				} else {
					record.RingingDuration = int(record.EndTime.Sub(*record.RingingTime) / time.Second)
				}
			}
			if record.AnswerTime != nil {
				record.TalkDuration = int(record.EndTime.Sub(*record.AnswerTime) / time.Second)
			}
		}
		go s.dealRTCPReport(callID)

		err := s.repository.CreateCall(ctx, record)
		if err != nil {
			logrus.WithError(err).Error("更新SIP呼叫记录失败")
		} else {
			// 从缓存中删除已保存的记录
			delete(s.callRecordCache, callID)
		}
	}
}

// 处理RTCP报告
func (s *SaveService) dealRTCPReport(callID string) {
	report := s.rtcpService.GetCallRTCPReportByCallID(callID)
	if report == nil {
		return
	}

	countLegs := len(report.Legs)
	if countLegs == 0 {
		return
	}

	rtcpReport := &entity.RtcpReport{}
	rtcpReportRaws := make([]*entity.RtcpReportRaw, 0)

	// 有两条或更多通信通道，设置A-leg和B-leg
	i := 0
	for _, leg := range report.Legs {
		leg.ProcessRTCPPackets()
		if i == 0 {
			rtcpReport.NodeIP = leg.NodeIP
			rtcpReport.SIPCallID = callID
			rtcpReport.SrcAddr = leg.SrcAddr
			rtcpReport.DstAddr = leg.DstAddr
			rtcpReport.CreateTime = time.Now()
			rtcpReport.TimestampMicro = time.Now().UnixMicro()

			// 创建A-leg RTCP报告
			rtcpReport.AlegMos = leg.Mos
			rtcpReport.AlegPacketLost = leg.PacketLost
			rtcpReport.AlegPacketCount = leg.PacketCount
			rtcpReport.AlegPacketLostRate = leg.PacketLostRate
			rtcpReport.AlegJitterAvg = leg.JitterAvg
			rtcpReport.AlegJitterMax = leg.JitterMax
			rtcpReport.AlegDelayAvg = leg.DelayAvg
			rtcpReport.AlegDelayMax = leg.DelayMax
		} else if i == 1 {
			rtcpReport.BlegMos = leg.Mos
			rtcpReport.BlegPacketLost = leg.PacketLost
			rtcpReport.BlegPacketCount = leg.PacketCount
			rtcpReport.BlegPacketLostRate = leg.PacketLostRate
			rtcpReport.BlegJitterAvg = leg.JitterAvg
			rtcpReport.BlegJitterMax = leg.JitterMax
			rtcpReport.BlegDelayAvg = leg.DelayAvg
			rtcpReport.BlegDelayMax = leg.DelayMax
		}

		for _, packet := range leg.RawPackets {
			rtcpReportRaws = append(rtcpReportRaws, &entity.RtcpReportRaw{
				NodeIP:     leg.NodeIP,
				SIPCallID:  callID,
				SrcAddr:    leg.SrcAddr,
				DstAddr:    leg.DstAddr,
				Raw:        packet.Raw,
				CreateTime: time.UnixMicro(packet.TimestampMicro),
			})
		}
		i++
	}

	ctx := context.Background()
	err := s.repository.CreateRtcpReport(ctx, rtcpReport)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"callID": callID,
			"report": rtcpReport,
		}).WithError(err).Error("保存RTCP报告失败")
	}

	if len(rtcpReportRaws) > 0 {
		err := s.repository.CreateRtcpReportRaws(ctx, rtcpReportRaws)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"callID": callID,
				"report": rtcpReportRaws,
			}).WithError(err).Error("保存RTCP报告失败")
		}
	}

	// 清理RTCP报告数据，避免内存泄漏
	s.rtcpService.ClearCallReport(callID)
}
