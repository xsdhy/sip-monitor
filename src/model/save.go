package model

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"sync"
	"time"

	"sip-monitor/src/entity"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 使用一个内存映射来缓存SIP呼叫记录
var callRecordCache = make(map[string]*entity.SIPRecordCall)
var cacheMutex = sync.RWMutex{}

// 定义用于测试的内部函数变量
var internalInsertOne = func(ctx context.Context, document interface{}, opts ...interface{}) (interface{}, error) {
	if CollectionRecord != nil {
		result, err := CollectionRecord.InsertOne(ctx, document)
		return result, err
	}
	return nil, nil
}

var internalUpdateOne = func(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) (interface{}, error) {
	if CollectionRecordCall != nil {
		ops := make([]*options.UpdateOptions, 0)
		for _, opt := range opts {
			if updateOpt, ok := opt.(*options.UpdateOptions); ok {
				ops = append(ops, updateOpt)
			}
		}
		result, err := CollectionRecordCall.UpdateOne(ctx, filter, update, ops...)
		return result, err
	}
	return nil, nil
}

// 定时将缓存刷新到数据库
func InitSaveToDBRunner() {
	// 启动周期性刷新缓存到数据库的任务
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				FlushCacheToDB()
			}
		}
	}()

	// 启动处理队列的任务
	go SaveToDBRunner()
}

// 将缓存刷新到数据库
func FlushCacheToDB() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	now := time.Now()
	count := 0

	// 遍历缓存中的记录
	for callID, record := range callRecordCache {
		// 检查记录是否已完成或超过一定时间未更新
		if record.EndTime != nil ||
			(record.CreateTime != nil && now.Sub(*record.CreateTime) > 15*time.Minute) {
			// 将记录保存到数据库
			ctx := context.Background()
			filter := bson.M{"sip_call_id": callID}
			opts := options.Update().SetUpsert(true)

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

			// 转换为bson.D
			data, err := bson.Marshal(record)
			if err != nil {
				logrus.WithError(err).Error("序列化SIP记录失败")
				continue
			}

			var updateDoc bson.D
			err = bson.Unmarshal(data, &updateDoc)
			if err != nil {
				logrus.WithError(err).Error("反序列化SIP记录失败")
				continue
			}

			// 使用内部函数进行更新，便于测试
			_, err = internalUpdateOne(ctx, filter, bson.D{{"$set", updateDoc}}, opts)
			if err != nil {
				logrus.WithError(err).Error("更新SIP呼叫记录失败")
			} else {
				count++
				// 从缓存中删除已保存的记录
				delete(callRecordCache, callID)
			}
		}
	}

	if count > 0 {
		logrus.WithField("count", count).Info("成功将缓存中的SIP呼叫记录写入数据库")
	}
}

func SaveToDBRunner() {
	for {
		select {
		case item := <-SaveToDBQueue:
			SaveOptimized(item)
		}
	}
}

// 优化的Save方法，减少数据库操作
func SaveOptimized(item entity.SIP) {
	if CollectionRecord == nil {
		return
	}

	// 始终需要在Record表中，新增一条记录
	go func() {
		// 将SIP转换为Record
		record := entity.Record{
			ID:         GetMd5(item.CallID, fmt.Sprintf("%v", item.TimestampMicroWithDate), item.NodeIP),
			NodeIP:     item.NodeIP,
			SIPCallID:  item.CallID,
			CreateTime: item.CreateAt,
			RawMsg:     *item.Raw,
			ViaNum:     item.ViaNum,
		}

		// 使用内部函数进行插入，便于测试
		_, err := internalInsertOne(context.TODO(), record)
		if err != nil {
			logrus.WithError(err).Error("保存SIP消息记录失败")
			return
		}
	}()

	// 忽略注册和通知消息
	if item.CSeqMethod == "REGISTER" || item.CSeqMethod == "NOTIFY" {
		return
	}

	// 使用内存缓存处理呼叫记录
	updateCallRecordInCache(item)
}

// 在内存缓存中更新呼叫记录
func updateCallRecordInCache(item entity.SIP) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	callID := item.CallID
	record, exists := callRecordCache[callID]

	// 如果记录不存在，创建一个新记录
	if !exists {
		// 对于新记录，只有INVITE方法才会创建
		if item.Title == "INVITE" && item.CSeqMethod == "INVITE" {
			record = &entity.SIPRecordCall{
				ID:        GetMd5(item.CallID, "", item.NodeIP),
				NodeIP:    item.NodeIP,
				SIPCallID: item.CallID,
			}

			// 设置基本字段
			if item.ViaNum == 1 {
				record.ToUser = item.ToUsername
				record.FromUser = item.FromUsername
				record.UserAgent = item.UserAgent
				record.SrcHost = item.SrcHost
				record.SrcPort = item.SrcPort
				record.SrcAddr = item.SrcAddr
				record.DstHost = item.DstHost
				record.DstPort = item.DstPort
				record.DstAddr = item.DstAddr
				record.TimestampMicro = int64(item.TimestampMicro)

				createTime := item.CreateAt
				record.CreateTime = &createTime
			}

			callRecordCache[callID] = record
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
				ringingTime := item.CreateAt
				record.RingingTime = &ringingTime
			}
			break
		case "200": // OK
			if (item.CSeqMethod == "ACK" || item.CSeqMethod == "INVITE") && record.AnswerTime == nil {
				answerTime := item.CreateAt
				record.AnswerTime = &answerTime
			} else if item.CSeqMethod == "BYE" && record.EndTime == nil {
				endTime := item.CreateAt
				record.EndTime = &endTime
				record.HangupCode = 200
				record.HangupCause = "Normal Clearing"
			}
			break
		case "CANCEL", "480", "487", "500": // Error or Cancel
			if record.EndTime == nil {
				endTime := item.CreateAt
				record.EndTime = &endTime
				// 设置挂断原因
				record.HangupCode = item.ResponseCode
				if item.Title == "CANCEL" {
					record.HangupCause = "Call Canceled"
				} else {
					record.HangupCause = item.ResponseDesc
				}
			}
			break
		}
	}

	// 如果记录已完成（已结束），立即写入数据库并从缓存移除
	if record.EndTime != nil {
		ctx := context.Background()
		filter := bson.M{"sip_call_id": callID}
		opts := options.Update().SetUpsert(true)

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

		// 转换为bson.D
		data, err := bson.Marshal(record)
		if err != nil {
			logrus.WithError(err).Error("序列化SIP记录失败")
			return
		}

		var updateDoc bson.D
		err = bson.Unmarshal(data, &updateDoc)
		if err != nil {
			logrus.WithError(err).Error("反序列化SIP记录失败")
			return
		}

		// 使用内部函数进行更新，便于测试
		_, err = internalUpdateOne(ctx, filter, bson.D{{"$set", updateDoc}}, opts)
		if err != nil {
			logrus.WithError(err).Error("更新SIP呼叫记录失败")
		} else {
			// 从缓存中删除已保存的记录
			delete(callRecordCache, callID)
		}
	}
}

func GetMd5(uuid string, content string, ip string) string {
	hash := md5.Sum([]byte(content))
	md5String := hex.EncodeToString(hash[:])

	return extractAlphanumeric(uuid) + md5String
}

func extractAlphanumeric(inputString string) string {
	// 定义正则表达式
	regex := regexp.MustCompile("[a-zA-Z0-9]")

	// 查找所有匹配项
	matches := regex.FindAllString(inputString, -1)

	// 构建字母和数字字符串
	result := ""
	for _, match := range matches {
		result += match
	}

	return result
}
