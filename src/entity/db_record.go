package entity

import (
	"time"
)

type Record struct {
	ID         string    `bson:"-" json:"id" gorm:"type:varchar(36);primaryKey;comment:'记录ID'"`
	UUID       string    `bson:"uuid" json:"uuid" gorm:"type:varchar(36);comment:'系统唯一ID'"`
	NodeID     string    `bson:"node_id" json:"node_id" gorm:"type:varchar(36);comment:'节点ID'"`
	NodeIP     string    `bson:"node_ip" json:"node_ip" gorm:"type:varchar(15);comment:'节点IP'"`
	CallID     string    `bson:"call_id" json:"call_id" gorm:"type:varchar(36);comment:'通话ID'"`
	Method     string    `bson:"method" json:"method" gorm:"type:varchar(36);comment:'方法'"`
	Src        string    `bson:"src" json:"src" gorm:"type:varchar(36);comment:'来源地址'"`
	Dst        string    `bson:"dst" json:"dst" gorm:"type:varchar(36);comment:'目的地址'"`
	CreateTime time.Time `bson:"create_time" json:"create_time" gorm:"type:timestamp;comment:'创建时间'"`
	Timestamp  uint64    `bson:"timestamp" json:"timestamp" gorm:"type:int;comment:'时间戳微秒部分'"`

	Body string `bson:"body" json:"body" gorm:"type:text;comment:'记录内容'"`
}
