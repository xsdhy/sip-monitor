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
	UUID string //唯一ID

	NodeIP string
	NodeID string

	Title       string // Method or Status
	SIPMethod   string
	SIPProtocol uint
	IsRequest   bool

	CallID          string
	RequestURL      string
	RequestUsername string
	RequestDomain   string

	ToUsername   string
	ToDomain     string
	FromUsername string
	FromDomain   string

	ResponseCode int
	ResponseDesc string

	CSeqNumber int
	CSeqMethod string
	UserAgent  string

	SrcHost        string
	SrcPort        int
	SrcAddr        string
	SrcCountryName string
	SrcCityName    string

	DstHost        string
	DstPort        int
	DstAddr        string
	DstCountryName string
	DstCityName    string

	CreateAt       time.Time
	TimestampMicro uint32
	Protocol       int

	Raw *string // raw sip message

	CreateTime time.Time
}
