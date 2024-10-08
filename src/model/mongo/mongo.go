package mongo

import (
	"context"
	"log/slog"

	"github.com/sirupsen/logrus"
	"sip-monitor/src/pkg/env"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MgWhereNotContain string = "^((?!%s).)*$"
	MgWhereContain    string = "^.*%s.*$"
)

type MgInfra struct {
	logger           *logrus.Logger
	MongoDB          *mongo.Database
	CollectionRecord *mongo.Collection

	CollectionRecordCall     *mongo.Collection
	CollectionRecordRegister *mongo.Collection
}

func NewMongoInfra(ctx context.Context, logger *logrus.Logger, dns string) (*MgInfra, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dns))
	if err != nil {
		slog.Error("MongoDBInit Error", slog.String("err:", err.Error()))
		return nil, err
	}
	m := &MgInfra{
		logger: logger,
	}
	m.MongoDB = client.Database(env.Conf.DBName)
	m.CollectionRecord = m.MongoDB.Collection("call_records")
	m.CollectionRecordCall = m.MongoDB.Collection("call_records_call")
	m.CollectionRecordRegister = m.MongoDB.Collection("call_records_register")

	index := mongo.IndexModel{
		Keys:    bson.M{"uuid": 1},
		Options: options.Index().SetUnique(true),
	}
	_, _ = m.CollectionRecordCall.Indexes().CreateOne(context.Background(), index)
	_, _ = m.CollectionRecordRegister.Indexes().CreateOne(context.Background(), index)
	return m, nil
}
