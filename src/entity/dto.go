package entity

import (
	"time"
)

type SearchParams struct {
	PageSize int64 `json:"page_size" form:"page_size"`
	Page     int64 `json:"page" form:"page"`

	NodeIP       string `json:"node_ip" form:"node_ip"`
	CallID       string `json:"call_id" form:"call_id"`
	UserAgent    string `json:"ua" form:"ua"`
	UserAgentOpr string `json:"ua_opr" form:"ua_opr"`

	SIPMethod       string `bson:"sip_method" json:"sip_method"`
	SIPMethodOpr    string `bson:"sip_method_opr" json:"sip_method_opr"`
	ResponseCode    int    `bson:"response_code" json:"response_code"`
	ResponseCodeOpr int    `bson:"response_code_opr" json:"response_code_opr"`

	BeginTime *time.Time `json:"begin_time" form:"begin_time" time_format:"2006-01-02 15:04:05" time_utc:"8"`
	EndTime   *time.Time `json:"end_time" form:"end_time" time_format:"2006-01-02 15:04:05" time_utc:"8"`

	FromUser          string `form:"from_user" json:"from_user"`
	FromUserOpr       string `form:"from_user_opr" json:"from_user_opr"`
	SrcHost           string `form:"src_host" json:"src_host"`
	SrcHostOpr        string `form:"src_host_opr" json:"src_host_opr"`
	SrcCountryName    string `form:"src_country_name" json:"src_country_name"`
	SrcCountryNameOpr string `form:"src_country_name_opr" json:"src_country_name_opr"`
	SrcCityName       string `form:"src_city_name" json:"src_city_name"`
	SrcCityNameOpr    string `form:"src_city_name_opr" json:"src_city_name_opr"`

	ToUser            string `form:"to_user" json:"to_user"`
	ToUserOpr         string `form:"to_user_opr" json:"to_user_opr"`
	DstHost           string `form:"dst_host" json:"dst_host"`
	DstHostOpr        string `form:"dst_host_opr" json:"dst_host_opr"`
	DstCountryName    string `form:"dst_country_name" json:"dst_country_name"`
	DstCountryNameOpr string `form:"dst_country_name_opr" json:"dst_country_name_opr"`
	DstCityName       string `form:"dst_city_name" json:"dst_city_name"`
	DstCityNameOpr    string `form:"dst_city_name_opr" json:"dst_city_name_opr"`
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
