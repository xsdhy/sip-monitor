package mysql

import (
	"context"
	"log/slog"
	"sip-monitor/src/entity"
	"sip-monitor/src/pkg/env"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type NoSqlInfra struct {
	db *gorm.DB
}

func NewNosqlInfra(ctx context.Context, dsn string) (*NoSqlInfra, error) {
	n := &NoSqlInfra{}
	var err error

	switch env.Conf.DBType {
	case "mysql":
		n.db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{})
	case "file":
		n.db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	case "memory":
		n.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	}

	_ = n.db.AutoMigrate(entity.Record{})
	_ = n.db.AutoMigrate(entity.SIPRecordRegister{})
	_ = n.db.AutoMigrate(entity.SIPRecordCall{})

	if err != nil {
		slog.Error("db err:", err.Error())
		return nil, err
	}

	return n, nil

}

func (n *NoSqlInfra) SaveMsg(sip *entity.SIP) {
	n.db.Table("record").Save(sip)
}

func (n *NoSqlInfra) SaveCall(sip *entity.SIPRecordCall) {
	n.db.Table("record_call").Save(sip)
}
