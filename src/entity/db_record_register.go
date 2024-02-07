package entity

import (
	"time"
)

type SIPRecordRegister struct {
	ID string `bson:"_id" json:"id"`

	NodeIP string `bson:"node_ip" json:"node_ip"`

	CreateTime time.Time `bson:"create_time" json:"create_time"`

	SIPCallID string `bson:"sip_call_id" json:"sip_call_id"`

	FromUser string `bson:"from_user" json:"from_user"`

	UserAgent string `bson:"user_agent" json:"user_agent"`

	RegisterTimes  int `bson:"register_times" json:"register_times"`   //注册次数
	FailuresTimes  int `bson:"failures_times" json:"failures_times"`   //失败次数:返回401,403的次数
	SuccessesTimes int `bson:"successes_times" json:"successes_times"` //成功次数:返回200的次数

	SrcHost        string `bson:"src_host" json:"src_host"`
	SrcPort        int    `bson:"src_port" json:"src_port"`
	SrcAddr        string `bson:"src_addr" json:"src_addr"`
	SrcCountryName string `bson:"src_country_name" json:"src_country_name"`
	SrcCityName    string `bson:"src_city_name" json:"src_city_name"`
}

type SIPRecordRegisterSaveDB struct {
	NodeIP string `bson:"node_ip" json:"node_ip"`

	CreateTime time.Time `bson:"create_time" json:"create_time"`

	SIPCallID string `bson:"sip_call_id" json:"sip_call_id"`

	FromUser string `bson:"from_user" json:"from_user"`

	UserAgent string `bson:"user_agent" json:"user_agent"`

	SrcHost        string `bson:"src_host" json:"src_host"`
	SrcPort        int    `bson:"src_port" json:"src_port"`
	SrcAddr        string `bson:"src_addr" json:"src_addr"`
	SrcCountryName string `bson:"src_country_name" json:"src_country_name"`
	SrcCityName    string `bson:"src_city_name" json:"src_city_name"`
}
