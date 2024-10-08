package callbuffer

import (
	"time"

	"sip-monitor/src/entity"
)

type RegBuffer struct {
	result entity.SIPRecordRegister
	stage  CallStage
}

func NewRegBuffer() *RegBuffer {
	return &RegBuffer{
		stage: StageInit,
	}
}

func (cb *RegBuffer) Add(item *entity.SIP) *entity.SIPRecordRegister {
	if item.CSeqMethod != "REGISTER" {
		return nil
	}
	if cb.stage >= StageCreate {
		//第二次
	} else {
		//第一次
		cb.stage = StageCreate
		cb.result.NodeID = item.NodeID
		cb.result.CreateTime = time.Now()
		cb.result.CallID = item.CallID

		cb.result.Username = item.FromUsername
		cb.result.UserAgent = item.UserAgent
		cb.result.SrcHost = item.SrcHost
		cb.result.SrcPort = item.SrcPort
		cb.result.SrcAddr = item.SrcAddr
		cb.result.SrcCountryName = item.SrcCountryName
		cb.result.SrcCityName = item.SrcCityName
	}
	switch item.Title {
	case "401", "403":
		//failures_times++
		cb.result.FailuresTimes++
		break
	case "200":
		//successes_times++
		cb.result.SuccessesTimes++
		break
	default:
		//register_times++
		cb.result.RegisterTimes++
	}

	return nil
}
