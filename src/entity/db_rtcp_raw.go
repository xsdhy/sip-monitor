package entity

import "time"

type RtcpReportRaw struct {
	ID     int64  `gorm:"primaryKey;column:id;type:bigint unsigned;autoIncrement:true" bson:"_id" json:"id"`
	NodeIP string `gorm:"column:node_ip;type:varchar(25);default:''" bson:"node_ip" json:"node_ip"`

	SIPCallID string `gorm:"column:sip_call_id;type:varchar(120);index;default:''" bson:"sip_call_id" json:"sip_call_id"`

	SrcAddr string `gorm:"column:src_addr;type:varchar(25);default:''" bson:"src_addr" json:"src_addr"` // Source address
	DstAddr string `gorm:"column:dst_addr;type:varchar(25);default:''" bson:"dst_addr" json:"dst_addr"` // Destination address

	Raw        string    `gorm:"column:raw;type:text" bson:"raw" json:"raw"`
	CreateTime time.Time `gorm:"column:create_time;index" bson:"create_time" json:"create_time"`
}

func (RtcpReportRaw) TableName() string {
	return "rtcp_report_raws"
}
