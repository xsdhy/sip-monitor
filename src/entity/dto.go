package entity

import (
	"time"
)

type SearchParams struct {
	PageSize int64  `json:"page_size" form:"page_size" query:"page_size"`
	Page     int64  `json:"page" form:"page" query:"page"`
	SortBy   string `json:"sort_by" form:"sort_by" query:"sort_by"`
	SortDesc bool   `json:"sort_desc" form:"sort_desc" query:"sort_desc"`

	SipCallID string `json:"sip_call_id" form:"sip_call_id" query:"sip_call_id"`
	SessionID string `json:"session_id" form:"session_id" query:"session_id"`

	BeginTime *time.Time `json:"begin_time" form:"begin_time" time_format:"2006-01-02 15:04:05" time_utc:"8" query:"begin_time"`
	EndTime   *time.Time `json:"end_time" form:"end_time" time_format:"2006-01-02 15:04:05" time_utc:"8" query:"end_time"`

	FromUser string `form:"from_user" json:"from_user" query:"from_user"`
	ToUser   string `form:"to_user" json:"to_user" query:"to_user"`

	SrcHost string `form:"src_host" json:"src_host" query:"src_host"`
	DstHost string `form:"dst_host" json:"dst_host" query:"dst_host"`

	HangupCode string `form:"hangup_code" json:"hangup_code" query:"hangup_code"`
}

type CleanSipRecordDTO struct {
	BeginTime *time.Time `json:"begin_time" form:"begin_time" time_format:"2006-01-02 15:04:05" time_utc:"8"`
	EndTime   *time.Time `json:"end_time" form:"end_time" time_format:"2006-01-02 15:04:05" time_utc:"8"`

	Method string `form:"method" json:"method"`
}

type AuthLogin struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type CallStatDTO struct {
	BeginTime *time.Time `json:"begin_time" form:"begin_time" time_format:"2006-01-02 15:04:05"`
	EndTime   *time.Time `json:"end_time" form:"end_time" time_format:"2006-01-02 15:04:05"`
}
