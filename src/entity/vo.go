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
	Records   []SIP `json:"records"`
	Relevants []SIP `json:"relevants"`
}
