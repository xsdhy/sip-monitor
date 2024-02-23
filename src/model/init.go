package model

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"sip-monitor/src/pkg/env"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database
var CollectionRecord *mongo.Collection

var CollectionRecordCall *mongo.Collection
var CollectionRecordRegister *mongo.Collection

func MongoDBInit() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if env.Conf.DSNURL == "" {
		slog.Info("MongoDBInit BY DBUser、DBPassword、DBAddr")
		env.Conf.DSNURL = fmt.Sprintf("mongodb://%s:%s@%s", env.Conf.DBUser, env.Conf.DBPassword, env.Conf.DBAddr)
	} else {
		slog.Info("MongoDBInit BY DSN_URL")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.Conf.DSNURL))
	if err != nil {
		slog.Error("MongoDBInit Error", slog.String("err:", err.Error()))
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
}
