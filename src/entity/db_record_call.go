package entity

import (
	"time"
)

type SIPRecordCall struct {
	ID string `bson:"_id" json:"id" gorm:"type:varchar(36);primaryKey;comment:'记录ID'"`

	UUID string `bson:"uuid" json:"uuid" gorm:"type:varchar(36);comment:'系统中的唯一ID'"`

	NodeID string `bson:"node_id" json:"node_id" gorm:"type:varchar(36);comment:'节点ID'"`
	NodeIP string `bson:"node_ip" json:"node_ip" gorm:"type:varchar(15);comment:'节点IP'"`
	CallID string `bson:"call_id" json:"call_id" gorm:"type:varchar(36);comment:'在node中通话唯一标识'"`

	FromUser string `bson:"from_user" json:"from_user" gorm:"type:varchar(50);comment:'发起用户'"`
	ToUser   string `bson:"to_user" json:"to_user" gorm:"type:varchar(50);comment:'接收用户'"`

	UserAgent string `bson:"user_agent" json:"user_agent" gorm:"type:varchar(256);comment:'用户代理'"`

	SrcHost        string `bson:"src_host" json:"src_host" gorm:"type:varchar(50);comment:'源主机'"`
	SrcPort        int    `bson:"src_port" json:"src_port" gorm:"type:int;comment:'源端口'"`
	SrcAddr        string `bson:"src_addr" json:"src_addr" gorm:"type:varchar(100);comment:'源地址'"`
	SrcCountryName string `bson:"src_country_name" json:"src_country_name" gorm:"type:varchar(50);comment:'源国家名称'"`
	SrcCityName    string `bson:"src_city_name" json:"src_city_name" gorm:"type:varchar(50);comment:'源城市名称'"`

	DstHost        string `bson:"dst_host" json:"dst_host" gorm:"type:varchar(50);comment:'目的主机'"`
	DstPort        int    `bson:"dst_port" json:"dst_port" gorm:"type:int;comment:'目的端口'"`
	DstAddr        string `bson:"dst_addr" json:"dst_addr" gorm:"type:varchar(100);comment:'目的地址'"`
	DstCountryName string `bson:"dst_country_name" json:"dst_country_name" gorm:"type:varchar(50);comment:'目的国家名称'"`
	DstCityName    string `bson:"dst_city_name" json:"dst_city_name" gorm:"type:varchar(50);comment:'目的城市名称'"`

	CreateTime  *time.Time `bson:"create_time" json:"create_time" gorm:"type:timestamp;comment:'创建时间'"`
	RingingTime *time.Time `bson:"ringing_time" json:"ringing_time" gorm:"type:timestamp;comment:'振铃时间'"`
	AnswerTime  *time.Time `bson:"answer_time" json:"answer_time" gorm:"type:timestamp;comment:'接听时间'"`
	EndTime     *time.Time `bson:"end_time" json:"end_time" gorm:"type:timestamp;comment:'结束时间'"`

	CallDuration    int `bson:"call_duration" json:"call_duration" gorm:"type:int;comment:'通话时长'"`
	RingingDuration int `bson:"ringing_duration" json:"ringing_duration" gorm:"type:int;comment:'振铃时长'"`
	TalkDuration    int `bson:"talk_duration" json:"talk_duration" gorm:"type:int;comment:'通话时间'"`

	HangupCode   string `bson:"hangup_code" json:"hangup_code" gorm:"type:varchar(10);comment:'挂断代码'"`
	HangupReason string `json:"hangup_reason" json:"hangup_reason" gorm:"type:varchar(256);comment:'挂断原因'"`
}
