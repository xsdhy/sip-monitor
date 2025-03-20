package entity

import (
	"time"
)

// SIPRecordCall represents a complete call record extracted from SIP signaling
type SIPRecordCall struct {
	// ID field - primary key
	ID int64 `gorm:"primaryKey;column:id;type:bigint unsigned;autoIncrement:true" bson:"_id" json:"id"`

	// NodeIP represents the IP of the node that collected the signal
	NodeIP string `gorm:"column:node_ip;type:varchar(25);index;default:''" bson:"node_ip" json:"node_ip"`

	// SIPCallID represents the unique call identifier from SIP protocol
	SIPCallID string `gorm:"column:sip_call_id;type:varchar(120);index;default:''" bson:"sip_call_id" json:"sip_call_id"`

	SessionID string `gorm:"column:session_id;type:varchar(120);index;default:''" bson:"session_id" json:"session_id"`

	// Call participants information
	ToUser    string `gorm:"column:to_user;type:varchar(120);index;default:''" bson:"to_user" json:"to_user"`
	FromUser  string `gorm:"column:from_user;type:varchar(120);index;default:''" bson:"from_user" json:"from_user"`
	UserAgent string `gorm:"column:user_agent;type:varchar(120);default:''" bson:"user_agent" json:"user_agent"` // User agent

	SrcAddr string `gorm:"column:src_addr;type:varchar(25);default:''" bson:"src_addr" json:"src_addr"` // Source address
	DstAddr string `gorm:"column:dst_addr;type:varchar(25);default:''" bson:"dst_addr" json:"dst_addr"` // Destination address

	// Timestamp in microseconds
	TimestampMicro int64 `gorm:"column:timestamp_micro;type:bigint unsigned;default:0" bson:"timestamp_micro" json:"timestamp_micro"`

	// Call timing information
	CreateTime  *time.Time `gorm:"column:create_time;index" bson:"create_time" json:"create_time"` // Creation time
	RingingTime *time.Time `gorm:"column:ringing_time" bson:"ringing_time" json:"ringing_time"`    // Ringing time
	AnswerTime  *time.Time `gorm:"column:answer_time" bson:"answer_time" json:"answer_time"`       // Answer time
	EndTime     *time.Time `gorm:"column:end_time" bson:"end_time" json:"end_time"`                // End time

	// Call duration measurements
	CallDuration    int `gorm:"column:call_duration;type:int unsigned;default:0" bson:"call_duration" json:"call_duration"`          // Total call duration
	RingingDuration int `gorm:"column:ringing_duration;type:int unsigned;default:0" bson:"ringing_duration" json:"ringing_duration"` // Ringing duration
	TalkDuration    int `gorm:"column:talk_duration;type:int unsigned;default:0" bson:"talk_duration" json:"talk_duration"`          // Talk duration

	// Call status information
	CallStatus  int    `gorm:"column:call_status;type:tinyint unsigned;default:0" bson:"call_status" json:"call_status"` // Call status
	HangupCode  int    `gorm:"column:hangup_code;type:tinyint unsigned;default:0" bson:"hangup_code" json:"hangup_code"` // Hangup code
	HangupCause string `gorm:"column:hangup_cause;type:varchar(120);default:''" bson:"hangup_cause" json:"hangup_cause"` // Hangup cause
}

// TableName specifies the database table name for GORM
func (SIPRecordCall) TableName() string {
	return "call_records_call"
}
