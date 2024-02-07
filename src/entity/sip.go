package entity

import (
	"time"
)

type SIP struct {
	NodeID string `json:"node_id"`
	NodeIP string `json:"node_ip"`

	Protocol int `json:"protocol"`

	CallID    string `json:"sip_call_id"`
	SessionID string `json:"session_id"`

	Title        string `json:"sip_method"` // Method or Status
	IsRequest    bool   `json:"is_request"`
	ResponseCode int    `json:"response_code"`
	ResponseDesc string `json:"response_desc"`

	CSeqNumber int    `json:"cseq_number"`
	CSeqMethod string `json:"cseq_method"`

	UserAgent string `json:"user_agent"`

	FromUser string `json:"from_user"`
	ToUser   string `json:"to_user"`

	SrcAddr string `json:"src_addr"`
	DstAddr string `json:"dst_addr"`

	CreateTime     time.Time `json:"create_time"`
	TimestampMicro int64     `json:"timestamp_micro"`

	Raw *string `json:"raw"` // raw sip message
}
