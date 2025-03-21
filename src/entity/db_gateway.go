package entity

import "time"

type Gateway struct {
	ID       int64      `gorm:"primaryKey;column:id;autoIncrement:true" bson:"_id" json:"id"`
	Name     string     `gorm:"column:name;type:varchar(120);default:''" bson:"name" json:"name"`
	Addr     string     `gorm:"column:addr;type:varchar(25);default:''" bson:"addr" json:"addr"`
	Remark   string     `gorm:"column:remark;type:varchar(255);default:''" bson:"remark" json:"remark"`
	CreateAt *time.Time `gorm:"column:create_at" bson:"create_at" json:"create_at"`
	UpdateAt *time.Time `gorm:"column:update_at" bson:"update_at" json:"update_at"`
}

func (Gateway) TableName() string {
	return "gateways"
}
