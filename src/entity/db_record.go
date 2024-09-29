package entity

import (
	"time"
)

type Record struct {
	ID         string    `bson:"_id" json:"id" gorm:"type:varchar(36);primaryKey;comment:'记录ID'"`
	UUID       string    `bson:"uuid" json:"uuid" gorm:"type:varchar(36);comment:'系统唯一ID'"`
	NodeID     string    `bson:"node_id" json:"node_id" gorm:"type:varchar(36);comment:'节点ID'"`
	NodeIP     string    `bson:"node_ip" json:"node_ip" gorm:"type:varchar(15);comment:'节点IP'"`
	CallID     string    `bson:"call_id" json:"call_id" gorm:"type:varchar(36);comment:'通话ID'"`
	CreateTime time.Time `bson:"create_time" json:"create_time" gorm:"type:timestamp;comment:'创建时间'"`
	Body       string    `bson:"body" json:"body" gorm:"type:text;comment:'记录内容'"`
}
