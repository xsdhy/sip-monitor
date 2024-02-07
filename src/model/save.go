package model

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"regexp"
	"sbc/src/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func upsertSIPRecordCallTime(sipCallID, timeType string, timeValue time.Time) {
	ctx := context.Background()
	filter := bson.D{{"sip_call_id", sipCallID}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", bson.D{
			{timeType, timeValue},
		}},
	}
	_, err := CollectionRecordCall.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		slog.Error("upsertSIPRecordCallTime save", err)
	}
}
func upsertSIPRecordCallInviteV3(record entity.Record, viaNum int) {
	ctx := context.Background()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"sip_call_id": record.SIPCallID}
	updateItems := bson.D{}

	if viaNum > 1 {
		updateItems = bson.D{
			{"dst_host", record.DstHost},
			{"dst_port", record.DstPort},
			{"dst_addr", record.DstAddr},
			{"dst_country_name", record.DstCountryName},
			{"dst_city_name", record.DstCityName},
		}
	} else {
		updateItems = bson.D{
			{"node_ip", record.NodeIP},
			{"sip_call_id", record.SIPCallID},

			{"to_user", record.ToUser},
			{"from_user", record.FromUser},

			{"user_agent", record.UserAgent},

			{"src_host", record.SrcHost},
			{"src_port", record.SrcPort},
			{"src_addr", record.SrcAddr},
			{"src_country_name", record.SrcCountryName},
			{"src_city_name", record.SrcCityName},

			{"dst_host", record.DstHost},
			{"dst_port", record.DstPort},
			{"dst_addr", record.DstAddr},
			{"dst_country_name", record.DstCountryName},
			{"dst_city_name", record.DstCityName},

			{"create_time", record.CreateTime},
		}
	}
	_, err := CollectionRecordCall.UpdateOne(ctx, filter, bson.D{{"$set", updateItems}}, opts)
	if err != nil {
		slog.Error("upsertSIPRecordCallInviteV3 save", err)
	}
}

func Save(item entity.Record, viaNum int) {
	if CollectionRecord == nil {
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
	case "INVITE", "BYE", "ACK", "CANCEL", "UPDATE":
		switch item.SIPMethod {
		case "INVITE":
			upsertSIPRecordCallInviteV3(item, viaNum)
			break
		case "180", "183":
			upsertSIPRecordCallTime(item.SIPCallID, "ringing_time", item.CreateTime)
			break
		case "200":
			if item.CSeqMethod == "ACK" || item.CSeqMethod == "INVITE" {
				upsertSIPRecordCallTime(item.SIPCallID, "answer_time", item.CreateTime)
			} else if item.CSeqMethod == "BYE" {
				upsertSIPRecordCallTime(item.SIPCallID, "end_time", item.CreateTime)
			}
			break
		case "CANCEL", "480":
			upsertSIPRecordCallTime(item.SIPCallID, "end_time", item.CreateTime)
			break
		case "100", "ACK", "BYE":
			break
		default:
			break
		}
	case "NOTIFY":
		return
	}

	//插入某一条数据
	_, err := CollectionRecord.InsertOne(ctx, item)
	if err != nil {
		slog.Error("Save Item Sip Message Error:", err.Error())
		return
	}
	slog.Info("Save Item", slog.String("msg", fmt.Sprintf("%s(%s) %s->%s", item.CSeqMethod, item.SIPCallID, item.FromUser+item.FromHost, item.ToUser+item.ToHost)))
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
