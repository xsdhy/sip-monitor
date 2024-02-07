package entity

import (
	"time"
)

type SearchParams struct {
	PageSize int64 `json:"page_size" form:"page_size"`
	Page     int64 `json:"page" form:"page"`

	NodeIP    string `json:"node_ip" form:"node_ip"`
	SipCallID string `json:"sip_call_id" form:"sip_call_id"`
	UserAgent string `json:"ua" form:"ua"`

	BeginTime *time.Time `json:"begin_time" form:"begin_time" time_format:"2006-01-02 15:04:05" time_utc:"8"`
	EndTime   *time.Time `json:"end_time" form:"end_time" time_format:"2006-01-02 15:04:05" time_utc:"8"`

	FromUser       string `form:"from_user" json:"from_user"`
	SrcHost        string `form:"src_host" json:"src_host"`
	SrcCountryName string `form:"src_country_name" json:"src_country_name"`
	SrcCityName    string `form:"src_city_name" json:"src_city_name"`

	ToUser         string `form:"to_user" json:"to_user"`
	DstHost        string `form:"dst_host" json:"dst_host"`
	DstCountryName string `form:"dst_country_name" json:"dst_country_name"`
	DstCityName    string `form:"dst_city_name" json:"dst_city_name"`
}

type CleanSipRecordDTO struct {
	BeginTime *time.Time `json:"begin_time" form:"begin_time" time_format:"2006-01-02 15:04:05" time_utc:"8"`
	EndTime   *time.Time `json:"end_time" form:"end_time" time_format:"2006-01-02 15:04:05" time_utc:"8"`

	Method string `form:"method" json:"method"`
}
