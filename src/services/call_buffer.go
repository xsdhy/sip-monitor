package services

import (
	"sip-monitor/src/entity"
)

type CallStage uint

const (
	StageInit CallStage = iota
	StageCreate
	StageRing
	StageAnswer
	StageEnd
)

type CallBuffer struct {
	result entity.SIPRecordCall
	stage  CallStage
}

func NewCallBuffer() *CallBuffer {
	return &CallBuffer{
		stage: StageInit,
	}
}

func (cb *CallBuffer) Add(item entity.Record, viaNum int) *entity.SIPRecordCall {
	switch item.CSeqMethod {
	case "REGISTER":
		break
	case "INVITE", "BYE", "ACK", "CANCEL", "UPDATE":
		switch item.SIPMethod {
		case "INVITE":
			if viaNum > 1 {
				//更新目的地址
				cb.result.DstHost = item.DstHost
				cb.result.DstPort = item.DstPort
				cb.result.DstAddr = item.DstAddr
				cb.result.DstCountryName = item.DstCountryName
				cb.result.DstCityName = item.DstCityName
			} else {
				//第一次创建
				cb.stage = StageCreate

				cb.result.NodeIP = item.NodeIP
				cb.result.SIPCallID = item.SIPCallID
				cb.result.ToUser = item.ToUser
				cb.result.FromUser = item.FromUser
				cb.result.UserAgent = item.UserAgent

				cb.result.SrcHost = item.SrcHost
				cb.result.SrcPort = item.SrcPort
				cb.result.SrcAddr = item.SrcAddr
				cb.result.SrcCountryName = item.SrcCountryName
				cb.result.SrcCityName = item.SrcCityName

				cb.result.DstHost = item.DstHost
				cb.result.DstPort = item.DstPort
				cb.result.DstAddr = item.DstAddr
				cb.result.DstCountryName = item.DstCountryName
				cb.result.DstCityName = item.DstCityName

				cb.result.CreateTime = &item.CreateTime
			}
			break
		case "180", "183":
			if cb.stage < StageRing {
				cb.stage = StageRing
			}
			cb.result.RingingTime = &item.CreateTime
			break
		case "200":
			if item.CSeqMethod == "ACK" || item.CSeqMethod == "INVITE" {
				if cb.stage < StageAnswer {
					cb.stage = StageAnswer
				}
				cb.result.AnswerTime = &item.CreateTime
			} else if item.CSeqMethod == "BYE" {
				cb.stage = StageEnd
				cb.result.EndTime = &item.CreateTime
				return &cb.result
			}
			break
		case "CANCEL", "480", "486", "487", "500":
			cb.stage = StageEnd
			cb.result.EndTime = &item.CreateTime
			return &cb.result
		case "100", "ACK", "BYE":
			break
		default:
			break
		}
	case "NOTIFY":
		return nil
	}
	return nil

}
