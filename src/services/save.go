package services

import (
	"context"
	"sync"
	"time"

	"sip-monitor/src/entity"
	"sip-monitor/src/model"

	"github.com/sirupsen/logrus"
)

type SaveService struct {
	logger          *logrus.Logger
	repository      model.Repository
	callRecordCache map[string]*entity.SIPRecordCall
	cacheMutex      sync.RWMutex
	SaveToDBQueue   chan entity.SIP
}

func NewSaveService(logger *logrus.Logger, repository model.Repository) *SaveService {
	s := &SaveService{
		logger:          logger,
		repository:      repository,
		callRecordCache: make(map[string]*entity.SIPRecordCall),
		cacheMutex:      sync.RWMutex{},
		SaveToDBQueue:   make(chan entity.SIP, 20000),
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
		// 检查记录是否已完成或超过一定时间未更新
		if record.EndTime != nil ||
			(record.CreateTime != nil && now.Sub(*record.CreateTime) > 15*time.Minute) {
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

			// 使用内部函数进行更新，便于测试
			err := s.repository.CreateSIPCallRecord(ctx, record)
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

// 优化的Save方法，减少数据库操作
func (s *SaveService) SaveOptimized(item entity.SIP) {
	// 忽略注册和通知消息
	if item.CSeqMethod == "REGISTER" || item.CSeqMethod == "NOTIFY" {
		return
	}

	// 始终需要在Record表中，新增一条记录
	go func() {
		// 将SIP转换为Record
		record := entity.Record{
			NodeIP:         item.NodeIP,
			SIPCallID:      item.CallID,
			SessionID:      item.SessionID,
			Method:         item.Title,
			SrcAddr:        item.SrcAddr,
			DstAddr:        item.DstAddr,
			CreateTime:     item.CreateTime,
			TimestampMicro: item.TimestampMicro,
			Raw:            *item.Raw,
		}
		err := s.repository.CreateRecord(context.TODO(), &record)
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
			record = &entity.SIPRecordCall{
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
			}
		case "200": // OK
			if (item.CSeqMethod == "ACK" || item.CSeqMethod == "INVITE") && record.AnswerTime == nil {
				answerTime := item.CreateTime
				record.AnswerTime = &answerTime
			} else if item.CSeqMethod == "BYE" && record.EndTime == nil {
				endTime := item.CreateTime
				record.EndTime = &endTime
				record.HangupCode = 200
				record.HangupCause = "Normal Clearing"
			}
		case "CANCEL", "480", "487", "500", "503", "504", "404", "403", "408", "413", "416", "486", "488", "513", "580": // Error or Cancel
			if record.EndTime == nil {
				endTime := item.CreateTime
				record.EndTime = &endTime
				// 设置挂断原因
				if record.HangupCode == 0 {
					record.HangupCode = item.ResponseCode
				}

				if item.Title == "CANCEL" {
					record.HangupCause = "Call Canceled"
				} else {
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
				record.RingingDuration = int(record.RingingTime.Sub(*record.CreateTime) / time.Second)
			}
			if record.AnswerTime != nil {
				record.TalkDuration = int(record.EndTime.Sub(*record.AnswerTime) / time.Second)
			}
		}

		err := s.repository.CreateSIPCallRecord(ctx, record)
		if err != nil {
			logrus.WithError(err).Error("更新SIP呼叫记录失败")
		} else {
			// 从缓存中删除已保存的记录
			delete(s.callRecordCache, callID)
		}
	}
}
