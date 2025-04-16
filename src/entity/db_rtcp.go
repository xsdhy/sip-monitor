package entity

import "time"

type RtcpReport struct {
	ID int64 `gorm:"primaryKey;column:id;type:bigint unsigned;autoIncrement:true" bson:"_id" json:"id"`

	NodeIP string `gorm:"column:node_ip;type:varchar(25);default:''" bson:"node_ip" json:"node_ip"`

	SIPCallID string `gorm:"column:sip_call_id;type:varchar(120);index;default:''" bson:"sip_call_id" json:"sip_call_id"`

	SrcAddr string `gorm:"column:src_addr;type:varchar(25);default:''" bson:"src_addr" json:"src_addr"` // Source address
	DstAddr string `gorm:"column:dst_addr;type:varchar(25);default:''" bson:"dst_addr" json:"dst_addr"` // Destination address

	AlegMos            float64 `gorm:"column:aleg_mos;type:float;default:0" bson:"aleg_mos" json:"aleg_mos"`                                        // 平均MOS
	AlegPacketLost     uint64  `gorm:"column:aleg_packet_lost;type:int unsigned;default:0" bson:"aleg_packet_lost" json:"aleg_packet_lost"`         // 总丢包数
	AlegPacketCount    uint64  `gorm:"column:aleg_packet_count;type:int unsigned;default:0" bson:"aleg_packet_count" json:"aleg_packet_count"`      // 总包数
	AlegPacketLostRate float64 `gorm:"column:aleg_packet_lost_rate;type:float;default:0" bson:"aleg_packet_lost_rate" json:"aleg_packet_lost_rate"` // 丢包率
	AlegJitterAvg      uint64  `gorm:"column:aleg_jitter_avg;type:int unsigned;default:0" bson:"aleg_jitter_avg" json:"aleg_jitter_avg"`            // 平均抖动
	AlegJitterMax      uint64  `gorm:"column:aleg_jitter_max;type:int unsigned;default:0" bson:"aleg_jitter_max" json:"aleg_jitter_max"`            // 抖动最大值
	AlegDelayAvg       uint64  `gorm:"column:aleg_delay_avg;type:int unsigned;default:0" bson:"aleg_delay_avg" json:"aleg_delay_avg"`               // 平均延迟
	AlegDelayMax       uint64  `gorm:"column:aleg_delay_max;type:int unsigned;default:0" bson:"aleg_delay_max" json:"aleg_delay_max"`               // 延迟最大值

	BlegMos            float64 `gorm:"column:bleg_mos;type:float;default:0" bson:"bleg_mos" json:"bleg_mos"`                                        // 平均MOS
	BlegPacketLost     uint64  `gorm:"column:bleg_packet_lost;type:int unsigned;default:0" bson:"bleg_packet_lost" json:"bleg_packet_lost"`         // 总丢包数
	BlegPacketCount    uint64  `gorm:"column:bleg_packet_count;type:int unsigned;default:0" bson:"bleg_packet_count" json:"bleg_packet_count"`      // 总包数
	BlegPacketLostRate float64 `gorm:"column:bleg_packet_lost_rate;type:float;default:0" bson:"bleg_packet_lost_rate" json:"bleg_packet_lost_rate"` // 丢包率
	BlegJitterAvg      uint64  `gorm:"column:bleg_jitter_avg;type:int unsigned;default:0" bson:"bleg_jitter_avg" json:"bleg_jitter_avg"`            // 平均抖动
	BlegJitterMax      uint64  `gorm:"column:bleg_jitter_max;type:int unsigned;default:0" bson:"bleg_jitter_max" json:"bleg_jitter_max"`            // 抖动最大值
	BlegDelayAvg       uint64  `gorm:"column:bleg_delay_avg;type:int unsigned;default:0" bson:"bleg_delay_avg" json:"bleg_delay_avg"`               // 平均延迟
	BlegDelayMax       uint64  `gorm:"column:bleg_delay_max;type:int unsigned;default:0" bson:"bleg_delay_max" json:"bleg_delay_max"`               // 延迟最大值

	// CreateTime represents when the record was created
	CreateTime     time.Time `gorm:"column:create_time;index" bson:"create_time" json:"create_time"`
	TimestampMicro int64     `gorm:"column:timestamp_micro;type:bigint unsigned;default:0" bson:"timestamp_micro" json:"timestamp_micro"`
}

func (RtcpReport) TableName() string {
	return "rtcp_report"
}
