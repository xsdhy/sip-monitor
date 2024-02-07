package entity

import "time"

type User struct {
	ID       int64     `gorm:"primaryKey;column:id;autoIncrement:true" bson:"_id" json:"id"`
	Nickname string    `gorm:"column:nickname;index" bson:"nickname" json:"nickname"`
	Username string    `gorm:"column:username;index" bson:"username" json:"username"`
	Password string    `gorm:"column:password" bson:"password" json:"password"`
	CreateAt time.Time `gorm:"column:create_at" bson:"create_at" json:"create_at"`
	UpdateAt time.Time `gorm:"column:update_at" bson:"update_at" json:"update_at"`
}

func (User) TableName() string {
	return "users"
}
