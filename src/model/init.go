package model

import (
	"context"
	"fmt"
	"log/slog"
	"sip-monitor/src/model/mongo"
	"sip-monitor/src/model/mysql"
	"sip-monitor/src/pkg/callbuffer"
	"time"

	"sip-monitor/src/entity"
	"sip-monitor/src/pkg/env"
)

var Infra DB

var SaveToDBQueue chan *entity.SIP

var CallBufferMap map[string]*callbuffer.CallBuffer

type DB interface {
	SaveMsg(sip *entity.SIP)
	SaveCall(sip *entity.SIPRecordCall)

	GetDetailsBySipCallID(ctx context.Context, searchParams entity.SearchParams) ([]entity.Record, error)
	GetRecordRegisterList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error)
	GetRecordCallList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error)
}

func DBInit() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	switch env.Conf.DBType {
	case "mongo":
		if env.Conf.DSNURL == "" {
			slog.Info("MongoDBInit BY DBUser、DBPassword、DBAddr")
			env.Conf.DSNURL = fmt.Sprintf("mongodb://%s:%s@%s", env.Conf.DBUser, env.Conf.DBPassword, env.Conf.DBAddr)
		} else {
			slog.Info("MongoDBInit BY DSN_URL")
		}
		Infra, err = mongo.NewMongoInfra(ctx, env.Conf.DSNURL)
	case "mysql":
	case "file":
	case "memory":
		Infra, err = mysql.NewNosqlInfra(ctx, env.Conf.DSNURL)
	default:
		slog.Info("no db")
	}
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Init DB success")
}
