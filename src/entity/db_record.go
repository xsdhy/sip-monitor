package entity

import "time"

type Record struct {
	ID int64 `gorm:"primaryKey;column:id;type:bigint unsigned;autoIncrement:true" bson:"_id" json:"id"`

	NodeIP string `gorm:"column:node_ip;type:varchar(25);default:''" bson:"node_ip" json:"node_ip"`

	SIPCallID string `gorm:"column:sip_call_id;type:varchar(120);index;default:''" bson:"sip_call_id" json:"sip_call_id"`

	Method       string `gorm:"column:method;type:varchar(10);default:''" bson:"method" json:"method"`
	ResponseDesc string `gorm:"column:response_desc;type:varchar(100);default:''" bson:"response_desc" json:"response_desc"`

	ToUser   string `gorm:"column:to_user;type:varchar(120);default:''" bson:"to_user" json:"to_user"`
	FromUser string `gorm:"column:from_user;type:varchar(120);default:''" bson:"from_user" json:"from_user"`

	SrcAddr string `gorm:"column:src_addr;type:varchar(25);default:''" bson:"src_addr" json:"src_addr"` // Source address
	DstAddr string `gorm:"column:dst_addr;type:varchar(25);default:''" bson:"dst_addr" json:"dst_addr"` // Destination address

	// CreateTime represents when the record was created
	CreateTime     time.Time `gorm:"column:create_time;index" bson:"create_time" json:"create_time"`
	TimestampMicro int64     `gorm:"column:timestamp_micro;type:bigint unsigned;default:0" bson:"timestamp_micro" json:"timestamp_micro"`
}

// TableName specifies the database table name for GORM
func (Record) TableName() string {
	return "call_records"
}
