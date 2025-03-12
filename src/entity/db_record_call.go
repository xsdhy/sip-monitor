package entity

import (
	"time"
)

type SIPRecordCall struct {
	ID string `bson:"_id" json:"id"`

	NodeIP string `bson:"node_ip" json:"node_ip"`

	SIPCallID string `bson:"sip_call_id" json:"sip_call_id"`

	ToUser   string `bson:"to_user" json:"to_user"`
	FromUser string `bson:"from_user" json:"from_user"`

	UserAgent string `bson:"user_agent" json:"user_agent"`

	SrcHost string `bson:"src_host" json:"src_host"`
	SrcPort int    `bson:"src_port" json:"src_port"`
	SrcAddr string `bson:"src_addr" json:"src_addr"`

	DstHost string `bson:"dst_host" json:"dst_host"`
	DstPort int    `bson:"dst_port" json:"dst_port"`
	DstAddr string `bson:"dst_addr" json:"dst_addr"`

	TimestampMicro int64 `bson:"timestamp_micro" json:"timestamp_micro"`

	CreateTime  *time.Time `bson:"create_time" json:"create_time"`
	RingingTime *time.Time `bson:"ringing_time" json:"ringing_time"`
	AnswerTime  *time.Time `bson:"answer_time" json:"answer_time"`
	EndTime     *time.Time `bson:"end_time" json:"end_time"`

	CallDuration    int `bson:"call_duration" json:"call_duration"`
	RingingDuration int `bson:"ringing_duration" json:"ringing_duration"`
	TalkDuration    int `bson:"talk_duration" json:"talk_duration"`

	CallStatus  int    `bson:"call_status" json:"call_status"`
	HangupCode  int    `bson:"hangup_code" json:"hangup_code"`
	HangupCause string `bson:"hangup_cause" json:"hangup_cause"`
}
