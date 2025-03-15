package entity

import (
	"time"
)

const (
	HeaderCallID = "Call-ID"
	HeaderFrom   = "From"
	HeaderTo     = "To"
	HeaderUA     = "User-Agent"
	HeaderCSeq   = "CSeq"
)

type SIP struct {
	NodeID   string `json:"node_id"`
	NodeIP   string `json:"node_ip"`
	NodeName string `json:"node_name"`

	Title           string `json:"sip_method"` // Method or Status
	IsRequest       bool   `json:"is_request"`
	ResponseCode    int    `json:"response_code"`
	ResponseDesc    string `json:"response_desc"`
	CallID          string `json:"sip_call_id"`
	RequestURL      string `json:"request_url"`
	RequestUsername string `json:"request_username"`
	RequestDomain   string `json:"request_domain"`
	ToUsername      string `json:"to_user"`
	ToDomain        string `json:"to_domain"`
	FromUsername    string `json:"from_user"`
	FromDomain      string `json:"from_domain"`
	CSeqNumber      int    `json:"cseq_number"`
	CSeqMethod      string `json:"cseq_method"`
	UserAgent       string `json:"user_agent"`

	SrcHost string `json:"src_host"`
	SrcPort int    `json:"src_port"`
	SrcAddr string `json:"src_addr"`

	DstHost string `json:"dst_host"`
	DstPort int    `json:"dst_port"`
	DstAddr string `json:"dst_addr"`

	CreateAt       time.Time `json:"create_time"`
	TimestampMicro int64     `json:"timestamp_micro"`
	Protocol       int       `json:"protocol"`
	UID            string    `json:"uid"`        // correlative id for AB call leg
	FSCallID       string    `json:"fs_call_id"` // freeswitch CallID
	Raw            *string   `json:"raw_msg"`    // raw sip message
}
