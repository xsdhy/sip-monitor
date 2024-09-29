package entity

import (
	"time"
)

type SIPRecordRegister struct {
	ID string `bson:"_id" json:"id" gorm:"type:varchar(36);primaryKey;comment:'记录ID'"`

	UUID   string `bson:"uuid" json:"uuid" gorm:"type:varchar(36);comment:'系统唯一ID'"`
	NodeID string `bson:"node_id" json:"node_id" gorm:"type:varchar(36);comment:'节点ID'"`
	CallID string `bson:"call_id" json:"call_id" gorm:"type:varchar(36);comment:'通话ID'"`

	CreateTime time.Time `bson:"create_time" json:"create_time" gorm:"type:timestamp;comment:'创建时间'"`

	Username string `bson:"username" json:"username" gorm:"type:varchar(50);comment:'用户名'"`

	UserAgent string `bson:"user_agent" json:"user_agent" gorm:"type:varchar(256);comment:'用户代理'"`

	RegisterTimes  int `bson:"register_times" json:"register_times" gorm:"type:int;comment:'注册次数'"`
	FailuresTimes  int `bson:"failures_times" json:"failures_times" gorm:"type:int;comment:'失败次数'"`
	SuccessesTimes int `bson:"successes_times" json:"successes_times" gorm:"type:int;comment:'成功次数'"`

	SrcHost        string `bson:"src_host" json:"src_host" gorm:"type:varchar(50);comment:'源主机'"`
	SrcPort        int    `bson:"src_port" json:"src_port" gorm:"type:int;comment:'源端口'"`
	SrcAddr        string `bson:"src_addr" json:"src_addr" gorm:"type:varchar(100);comment:'源地址'"`
	SrcCountryName string `bson:"src_country_name" json:"src_country_name" gorm:"type:varchar(50);comment:'源国家名称'"`
	SrcCityName    string `bson:"src_city_name" json:"src_city_name" gorm:"type:varchar(50);comment:'源城市名称'"`
}
