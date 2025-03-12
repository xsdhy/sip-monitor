package model

import (
	"context"
	"fmt"
	"sip-monitor/src/entity"
	"time"

	"sip-monitor/src/pkg/env"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database
var CollectionRecord *mongo.Collection

var CollectionRecordCall *mongo.Collection
var CollectionRecordRegister *mongo.Collection

var SaveToDBQueue chan entity.SIP

func MongoDBInit() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if env.Conf.DSNURL == "" {
		logrus.Info("MongoDBInit BY DBUser、DBPassword、DBAddr")
		env.Conf.DSNURL = fmt.Sprintf("mongodb://%s:%s@%s", env.Conf.DBUser, env.Conf.DBPassword, env.Conf.DBAddr)
	} else {
		logrus.Info("MongoDBInit BY DSN_URL")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.Conf.DSNURL))
	if err != nil {
		logrus.WithError(err).Error("MongoDBInit Error")
		return
	}

	MongoDB = client.Database(env.Conf.DBName)
	CollectionRecord = MongoDB.Collection("call_records")
	CollectionRecordCall = MongoDB.Collection("call_records_call")
	CollectionRecordRegister = MongoDB.Collection("call_records_register")

	index := mongo.IndexModel{
		Keys:    bson.M{"sip_call_id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, _ = CollectionRecordCall.Indexes().CreateOne(context.Background(), index)
	_, _ = CollectionRecordRegister.Indexes().CreateOne(context.Background(), index)

	// 初始化缓存和保存机制
	InitSaveToDBRunner()
}
