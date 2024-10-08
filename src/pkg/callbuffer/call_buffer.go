package callbuffer

import (
	"time"

	"github.com/sirupsen/logrus"
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

type saveCall func(item *entity.SIPRecordCall)


type CallBuffer struct {
	logger      *logrus.Logger
	result      entity.SIPRecordCall
	stage       CallStage
	stageChange chan CallStage
	timer       *time.Timer
	saveCall    saveCall
}

func NewCallBuffer(logger *logrus.Logger, saveCall saveCall) *CallBuffer {
	cb := &CallBuffer{
		logger:      logger,
		stage:       StageInit,
		saveCall:    saveCall,
		stageChange: make(chan CallStage),
		timer:       time.NewTimer(70 * time.Second),
	}
	go cb.listener()
	return cb
}
func (cb *CallBuffer) listener() {
	for {
		select {
		case <-cb.timer.C:
			cb.logger.WithField("uuid", cb.result.UUID).Debug("time ok")
			cb.saveCall(&cb.result)
			return
		case stage := <-cb.stageChange:
			cb.stage = stage
			switch stage {
			case StageAnswer:
				cb.logger.WithField("uuid", cb.result.UUID).Debug("answer")
				cb.timer.Reset(30 * time.Minute)
			case StageEnd:
				cb.logger.WithField("uuid", cb.result.UUID).Debug("end")
				cb.timer.Reset(10 * time.Second)
			default:

			}
		}
	}
}

func (cb *CallBuffer) Add(item *entity.SIP) {
	switch item.CSeqMethod {
	case "REGISTER":
		break
	case "INVITE", "BYE", "ACK", "CANCEL", "UPDATE":
		switch item.Title {
		case "INVITE":
			if cb.stage >= StageCreate {
				//更新目的地址
				cb.result.DstHost = item.DstHost
				cb.result.DstPort = item.DstPort
				cb.result.DstCountryName = item.DstCountryName
				cb.result.DstCityName = item.DstCityName
			} else {
				//第一次创建
				cb.stageChange <- StageCreate

				cb.result.UUID = item.UUID
				cb.result.NodeID = item.NodeID
				cb.result.NodeIP = item.NodeIP

				cb.result.CallID = item.CallID
				cb.result.ToUser = item.ToUsername
				cb.result.FromUser = item.FromUsername
				cb.result.UserAgent = item.UserAgent

				cb.result.SrcHost = item.SrcHost
				cb.result.SrcPort = item.SrcPort
				cb.result.SrcCountryName = item.SrcCountryName
				cb.result.SrcCityName = item.SrcCityName

				cb.result.DstHost = item.DstHost
				cb.result.DstPort = item.DstPort
				cb.result.DstCountryName = item.DstCountryName
				cb.result.DstCityName = item.DstCityName

				cb.result.CreateTime = &item.CreateTime
			}
			break
		case "180", "183":
			if cb.stage < StageRing {
				cb.stageChange <- StageRing
			}
			cb.result.RingingTime = &item.CreateTime
			break
		case "200":
			if item.CSeqMethod == "ACK" || item.CSeqMethod == "INVITE" {
				if cb.stage < StageAnswer {
					cb.stageChange <- StageAnswer
				}
				cb.result.AnswerTime = &item.CreateTime
			} else if item.CSeqMethod == "BYE" {
				cb.stageChange <- StageEnd
				cb.result.EndTime = &item.CreateTime
			}
			cb.logger.WithField("uuid", cb.result.UUID).Debug("StageEnd by 200")
			break
		case "CANCEL", "480", "486", "487", "500":
			cb.stageChange <- StageEnd
			cb.result.EndTime = &item.CreateTime
			cb.logger.WithField("uuid", cb.result.UUID).Debug("StageEnd by other")
		case "100", "ACK", "BYE":
			break
		default:
			break
		}
		cb.duration()
	case "NOTIFY":
		return
	}
	return

}
func (cb *CallBuffer) duration() {
	if cb.result.CreateTime == nil {
		return
	}
	if cb.result.EndTime != nil {
		cb.result.CallDuration = int(cb.result.EndTime.Sub(*cb.result.CreateTime).Seconds())
		if cb.result.AnswerTime != nil {
			if cb.result.RingingTime != nil {
				cb.result.TalkDuration = int(cb.result.EndTime.Sub(*cb.result.RingingTime).Seconds())
			} else {
				cb.result.TalkDuration = cb.result.CallDuration
			}
		}
	}
	if cb.result.AnswerTime != nil {
		if cb.result.RingingTime != nil {
			cb.result.RingingDuration = int(cb.result.AnswerTime.Sub(*cb.result.RingingTime).Seconds())
		} else {
			cb.result.RingingDuration = int(cb.result.AnswerTime.Sub(*cb.result.CreateTime).Seconds())
		}
	}

}
