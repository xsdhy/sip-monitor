package entity

import (
	"time"
)

type Record struct {
	ID string `bson:"_id" json:"id"`

	NodeIP string `bson:"node_ip" json:"node_ip"`

	CreateTime time.Time `bson:"create_time" json:"create_time"`

	SIPCallID string `bson:"sip_call_id" json:"sip_call_id"`
	SIPMethod string `bson:"sip_method" json:"sip_method"`

	FsCallID string `bson:"fs_call_id" json:"fs_call_id"`

	LegUid   string `bson:"leg_uid" json:"leg_uid"`
	ToUser   string `bson:"to_user" json:"to_user"`
	FromUser string `bson:"from_user" json:"from_user"`

	ResponseCode int    `bson:"response_code" json:"response_code"`
	ResponseDesc string `bson:"response_desc" json:"response_desc"`
	CSeqMethod   string `bson:"cseq_method" json:"cseq_method"`
	CSeqNumber   int    `bson:"cseq_number" json:"cseq_number"`

	FromHost string `bson:"from_host" json:"from_host"`
	ToHost   string `bson:"to_host" json:"to_host"`

	SIPProtocol uint   `bson:"sip_protocol" json:"sip_protocol"`
	UserAgent   string `bson:"user_agent" json:"user_agent"`

	SrcHost string `bson:"src_host" json:"src_host"`
	SrcPort int    `bson:"src_port" json:"src_port"`
	SrcAddr string `bson:"src_addr" json:"src_addr"`

	DstHost string `bson:"dst_host" json:"dst_host"`
	DstPort int    `bson:"dst_port" json:"dst_port"`
	DstAddr string `bson:"dst_addr" json:"dst_addr"`

	TimestampMicro int64 `bson:"timestamp_micro" json:"timestamp_micro"`

	RawMsg string `bson:"raw_msg" json:"raw_msg"`

	ViaNum int `bson:"-" json:"-"` //临时使用
}
