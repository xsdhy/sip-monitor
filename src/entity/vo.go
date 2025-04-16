package entity

import (
	"time"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Time time.Time   `json:"time"`
}

type ResponseItems struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta"`
	Time time.Time   `json:"time"`
}

type Meta struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
}

type CallDetailsVO struct {
	Records     []Record         `json:"records"`
	Relevants   []Record         `json:"relevants"`
	RtcpReport  *RtcpReport      `json:"rtcp_report"`
	RTCPPackets []*RtcpReportRaw `json:"rtcp_packets"`
}

type CallStatVO struct {
	IP                 string `json:"ip"`
	Gateway            string `json:"gateway"`
	Total              int    `json:"total"`
	Answered           int    `json:"answered"`
	HangupCode0Count   int    `json:"hangup_code_0_count" gorm:"column:hangup_code_0_count"`
	HangupCode1XXCount int    `json:"hangup_code_1xx_count" gorm:"column:hangup_code_1xx_count"`
	HangupCode2XXCount int    `json:"hangup_code_2xx_count" gorm:"column:hangup_code_2xx_count"`
	HangupCode3XXCount int    `json:"hangup_code_3xx_count" gorm:"column:hangup_code_3xx_count"`
	HangupCode4XXCount int    `json:"hangup_code_4xx_count" gorm:"column:hangup_code_4xx_count"`
	HangupCode5XXCount int    `json:"hangup_code_5xx_count" gorm:"column:hangup_code_5xx_count"`
}
