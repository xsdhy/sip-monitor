package model

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"regexp"

	"sip-monitor/src/entity"
	"sip-monitor/src/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveToDBRunner() {
	for {
		select {
		case item := <-SaveToDBQueue:
			if item.CSeqMethod != "REGISTER" {
				saveItem(item)

				buffer, ok := CallBufferMap[item.SIPCallID]
				if !ok {
					CallBufferMap[item.SIPCallID] = services.NewCallBuffer()
				}
				result := buffer.Add(item, item.ViaNum)
				if result != nil {
					saveCall(result)
				}
			} else {
				saveRegister(item)
			}
		}
	}
}

func saveCall(item *entity.SIPRecordCall) {
	if CollectionRecordCall == nil || item == nil {
		return
	}
	//插入某一条数据
	_, err := CollectionRecordCall.InsertOne(context.Background(), item)
	if err != nil {
		slog.Error("Save Item call Error:", err.Error())
		return
	}
	slog.Debug("Save Item call", slog.String("msg", fmt.Sprintf("%s(%s) %s->%s", "inteve", item.SIPCallID, item.FromUser+item.SrcHost, item.ToUser+item.DstHost)))
}

func saveRegister(item entity.Record) {
	if CollectionRecordRegister == nil {
		return
	}

	ctx := context.TODO()
	filterSipCallID := bson.D{{"sip_call_id", item.SIPCallID}}
	opts := options.Update().SetUpsert(true)
	switch item.CSeqMethod {
	case "REGISTER":
		update := bson.D{}
		var updateFields bson.D
		conv, _ := bson.Marshal(entity.SIPRecordRegisterSaveDB{
			NodeIP:     item.NodeIP,
			CreateTime: item.CreateTime,
			SIPCallID:  item.SIPCallID,
			FromUser:   item.FromUser,
			UserAgent:  item.UserAgent,

			SrcHost:        item.SrcHost,
			SrcPort:        item.SrcPort,
			SrcAddr:        item.SrcAddr,
			SrcCountryName: item.SrcCountryName,
			SrcCityName:    item.SrcCityName,
		})
		_ = bson.Unmarshal(conv, &updateFields)

		switch item.SIPMethod {
		case "401", "403":
			update = bson.D{
				{"$inc", bson.D{{"failures_times", 1}}},
				{"$set", bson.D{{"sip_call_id", item.SIPCallID}}},
			}
			break
		case "200":
			update = bson.D{
				{"$inc", bson.D{{"successes_times", 1}}},
				{"$set", bson.D{{"sip_call_id", item.SIPCallID}}},
			}
			break
		default:
			update = bson.D{
				{"$inc", bson.D{{"register_times", 1}}},
				{"$set", updateFields},
			}
		}
		_, err := CollectionRecordRegister.UpdateOne(ctx, filterSipCallID, update, opts)
		if err != nil {
			slog.Error("register save", err)
		}
		break
	case "NOTIFY":
		return
	}

}

func saveItem(item entity.Record) {
	if CollectionRecord == nil {
		return
	}
	ctx := context.TODO()
	_, err := CollectionRecord.InsertOne(ctx, item)
	if err != nil {
		slog.Error("Save Item Sip Message Error:", err.Error())
		return
	}
	slog.Debug("Save Item", slog.String("msg", fmt.Sprintf("%s(%s) %s->%s", item.CSeqMethod, item.SIPCallID, item.FromUser+item.FromHost, item.ToUser+item.ToHost)))
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
