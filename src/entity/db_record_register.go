package entity

import (
	"time"
)

// SIPRecordRegister represents a SIP registration record
type SIPRecordRegister struct {
	// ID field - primary key
	ID int64 `gorm:"primaryKey;column:id;autoIncrement:true" bson:"_id" json:"id"`

	// NodeIP represents the IP of the node that collected the signal
	NodeIP string `gorm:"column:node_ip" bson:"node_ip" json:"node_ip"`

	// CreateTime represents when the registration record was created
	CreateTime *time.Time `gorm:"column:create_time;index" bson:"create_time" json:"create_time"`

	// SIPCallID represents the unique call identifier from SIP protocol
	SIPCallID string `gorm:"column:sip_call_id;index" bson:"sip_call_id" json:"sip_call_id"`

	// FromUser represents the user who is registering
	FromUser string `gorm:"column:from_user;index" bson:"from_user" json:"from_user"`

	// UserAgent represents the client software
	UserAgent string `gorm:"column:user_agent" bson:"user_agent" json:"user_agent"`

	// Registration statistics
	RegisterTimes  int `gorm:"column:register_times" bson:"register_times" json:"register_times"`    // Number of registration attempts
	FailuresTimes  int `gorm:"column:failures_times" bson:"failures_times" json:"failures_times"`    // Number of failures (401, 403 responses)
	SuccessesTimes int `gorm:"column:successes_times" bson:"successes_times" json:"successes_times"` // Number of successful registrations (200 responses)

	// Source information
	SrcAddr string `gorm:"column:src_addr" bson:"src_addr" json:"src_addr"`
}

// TableName specifies the database table name for GORM
func (SIPRecordRegister) TableName() string {
	return "call_records_register"
}
