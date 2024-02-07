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
	UID            string  // correlative id for AB call leg
	FSCallID       string  // freeswitch CallID
	Raw            *string // raw sip message
}
