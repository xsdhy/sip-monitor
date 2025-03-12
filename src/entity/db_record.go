package entity

import "time"

type Record struct {
	ID string `bson:"_id" json:"id"`

	NodeIP string `bson:"node_ip" json:"node_ip"`

	SIPCallID string `bson:"sip_call_id" json:"sip_call_id"`

	CreateTime time.Time `bson:"create_time" json:"create_time"`

	RawMsg string `bson:"raw_msg" json:"raw_msg"`

	ViaNum int `bson:"-" json:"-"` //临时使用
}
