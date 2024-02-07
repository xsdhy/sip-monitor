package entity

import "time"

// Record represents the original SIP signaling record
type Record struct {
	// ID field - primary key
	ID int64 `gorm:"primaryKey;column:id;autoIncrement:true" bson:"_id" json:"id"`

	// NodeIP represents the IP of the node that collected the signal
	NodeIP string `gorm:"column:node_ip;index" bson:"node_ip" json:"node_ip"`

	// SIPCallID represents the unique call identifier from SIP protocol
	SIPCallID string `gorm:"column:sip_call_id;index" bson:"sip_call_id" json:"sip_call_id"`

	SessionID string `gorm:"column:session_id;index" bson:"session_id" json:"session_id"`

	Method string `gorm:"column:method;index" bson:"method" json:"method"`

	SrcAddr string `gorm:"column:src_addr" bson:"src_addr" json:"src_addr"` // Source address
	DstAddr string `gorm:"column:dst_addr" bson:"dst_addr" json:"dst_addr"` // Destination address

	// CreateTime represents when the record was created
	CreateTime     time.Time `gorm:"column:create_time;index" bson:"create_time" json:"create_time"`
	TimestampMicro int64     `gorm:"column:timestamp_micro" bson:"timestamp_micro" json:"timestamp_micro"`

	Raw string `gorm:"column:raw;type:text" bson:"raw" json:"raw"`
}

// TableName specifies the database table name for GORM
func (Record) TableName() string {
	return "call_records"
}
