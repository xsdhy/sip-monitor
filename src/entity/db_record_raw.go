package entity

import "time"

type RecordRaw struct {
	ID         int64     `gorm:"primaryKey;type:bigint unsigned;autoIncrement:false" bson:"_id" json:"id"`
	Raw        string    `gorm:"column:raw;type:text" bson:"raw" json:"raw"`
	CreateTime time.Time `gorm:"column:create_time;index" bson:"create_time" json:"create_time"`
}

func (RecordRaw) TableName() string {
	return "call_record_raws"
}
