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
	NodeID   string
	NodeIP   string
	NodeName string

	Title           string // Method or Status
	IsRequest       bool
	ResponseCode    int
	ResponseDesc    string
	CallID          string
	RequestURL      string
	RequestUsername string
	RequestDomain   string
	ToUsername      string
	ToDomain        string
	FromUsername    string
	FromDomain      string
	CSeqNumber      int
	CSeqMethod      string
	UserAgent       string

	SrcHost string
	SrcPort int
	SrcAddr string

	DstHost string
	DstPort int
	DstAddr string

	CreateAt               time.Time
	TimestampMicro         uint32
	TimestampMicroWithDate int64
	Protocol               int
	UID                    string  // correlative id for AB call leg
	FSCallID               string  // freeswitch CallID
	Raw                    *string // raw sip message

	ViaNum int
}
