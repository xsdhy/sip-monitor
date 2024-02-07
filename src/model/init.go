package model

import (
	"context"
	"log/slog"
	"time"

	"sbc/src/pkg/env"

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

	slog.Info("MongoDBInit", slog.Bool("dsn", len(env.Conf.DSNURL) > 0))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.Conf.DSNURL))
	if err != nil {
		slog.Error("MongoDBInit Error", err.Error())
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
