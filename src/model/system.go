package model

import (
	"context"
	"time"

	"sip-monitor/src/entity"

	"go.mongodb.org/mongo-driver/bson"
)

func CleanSipRecord(request entity.CleanSipRecordDTO) (int64, error) {
	if request.EndTime == nil {
		return 0, nil
	}
	//filter := bson.M{}
	//
	//timeFilter := bson.M{}
	//if request.BeginTime != nil {
	//	timeFilter["$gte"] = request.BeginTime
	//}
	//timeFilter["$lte"] = request.EndTime
	//
	//filter["create_time"] = timeFilter
	//
	//if request.Method != "" {
	//	filter["cseq_method"] = request.Method
	//}
	//deleteResult, err := CollectionRecord.DeleteMany(context.Background(), filter, nil)
	//if err != nil {
	//	return 0, err
	//}
	//return deleteResult.DeletedCount, nil
	return 0, nil
}

func CleanSipALL(endTime *time.Time) (int64, error) {
	if endTime == nil {
		return 0, nil
	}
	filter := bson.M{}
	timeFilter := bson.M{}
	timeFilter["$lte"] = endTime
	filter["create_time"] = timeFilter

	var deleteTotal int64

	//record, err := CollectionRecord.DeleteMany(context.Background(), filter, nil)
	//if err == nil {
	//	deleteTotal = deleteTotal + record.DeletedCount
	//}
	//call, err := CollectionRecordCall.DeleteMany(context.Background(), filter, nil)
	//if err == nil {
	//	deleteTotal = deleteTotal + call.DeletedCount
	//}
	//reg, err := CollectionRecordRegister.DeleteMany(context.Background(), filter, nil)
	//if err == nil {
	//	deleteTotal = deleteTotal + reg.DeletedCount
	//}
	return deleteTotal, nil
}

func DbStats(ctx context.Context, tableName string) *entity.MongoDBStatsVO {
	//result := MongoDB.RunCommand(ctx, bson.D{{"collStats", tableName}})
	//if result.Err() != nil {
	//	return nil
	//}
	var stats entity.MongoDBStatsVO
	//if err := result.Decode(&stats); err != nil {
	//	return nil
	//}
	stats.Name = tableName
	return &stats
}
