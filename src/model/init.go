package model

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"sip-monitor/src/entity"
	"sip-monitor/src/model/mongo"
	"sip-monitor/src/model/mysql"
	"sip-monitor/src/pkg/env"
)

type DB interface {
	SaveMsg(sip *entity.SIP)
	SaveCall(sip *entity.SIPRecordCall)

	GetDetailsBySipCallID(ctx context.Context, searchParams entity.SearchParams) ([]entity.Record, error)
	GetRecordRegisterList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error)
	GetRecordCallList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error)
}

func DBInit(logger *logrus.Logger) (db DB, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch env.Conf.DBType {
	case "mongo":
		if env.Conf.DSNURL == "" {
			logger.Info("MongoDBInit BY DBUser、DBPassword、DBAddr")
			env.Conf.DSNURL = fmt.Sprintf("mongodb://%s:%s@%s", env.Conf.DBUser, env.Conf.DBPassword, env.Conf.DBAddr)
		} else {
			logger.Info("MongoDBInit BY DSN_URL")
		}
		db, err = mongo.NewMongoInfra(ctx, logger, env.Conf.DSNURL)
	case "mysql", "file", "memory":
		db, err = mysql.NewNosqlInfra(ctx, env.Conf.DSNURL)
	default:
		logger.Info("no db")
	}
	if err != nil {
		logger.WithError(err).Error("init db error")
	}
	logger.Info("Init DB success")
	return
}
