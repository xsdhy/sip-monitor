package model

import (
	"context"
	"fmt"

	"sbc/src/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSearchFilter(sp entity.SearchParams) bson.M {
	filter := bson.M{}
	if sp.BeginTime != nil && sp.EndTime != nil {
		filter["create_time"] = bson.M{
			"$gte": sp.BeginTime,
			"$lte": sp.EndTime,
		}
	}

	if sp.NodeIP != "" {
		filter["node_ip"] = sp.NodeIP
	}
	if sp.SipCallID != "" {
		filter["sip_call_id"] = sp.SipCallID
	}
	if sp.UserAgent != "" {
		filter["user_agent"] = bson.M{"$regex": fmt.Sprintf("^%s", sp.UserAgent)}
	}

	if sp.FromUser != "" {
		filter["from_user"] = bson.M{"$regex": fmt.Sprintf("^%s", sp.FromUser)}
	}
	if sp.SrcHost != "" {
		filter["scr_host"] = sp.SrcHost
	}
	if sp.ToUser != "" {
		filter["to_user"] = bson.M{"$regex": fmt.Sprintf("^%s", sp.ToUser)}
	}
	if sp.DstHost != "" {
		filter["to_host"] = sp.DstHost
	}

	return filter
}

func GetDetailsBySipCallID(ctx context.Context, searchParams entity.SearchParams) ([]entity.Record, error) {
	opt := options.Find().SetSort(bson.D{
		{"timestamp_micro", 1},
	})
	filter := GetSearchFilter(searchParams)

	cursor, err := CollectionRecord.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	records := make([]entity.Record, 0)
	// 处理查询结果
	for cursor.Next(ctx) {
		var result entity.Record
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		records = append(records, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func GetRecordList(ctx context.Context, searchParams entity.SearchParams) ([]entity.Record, *entity.Meta, error) {
	if searchParams.PageSize <= 0 {
		searchParams.PageSize = 10
	}
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}

	opt := options.Find().SetLimit(searchParams.PageSize).
		SetSkip(searchParams.PageSize * (searchParams.Page - 1)).
		SetSort(bson.D{{"create_time", -1}})

	filter := GetSearchFilter(searchParams)

	documentsCount, err := CollectionRecord.CountDocuments(ctx, filter, nil)
	if err != nil {
		return nil, nil, err
	}
	meta := entity.Meta{
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    documentsCount,
	}

	cursor, err := CollectionRecord.Find(ctx, filter, opt)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)
	records := make([]entity.Record, 0, searchParams.PageSize)
	// 处理查询结果
	for cursor.Next(ctx) {
		var result entity.Record
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		records = append(records, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, nil, err
	}
	return records, &meta, nil
}

func GetRecordRegisterList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error) {
	if searchParams.PageSize <= 0 {
		searchParams.PageSize = 10
	}
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}

	opt := options.Find().SetLimit(searchParams.PageSize).
		SetSkip(searchParams.PageSize * (searchParams.Page - 1)).
		SetSort(bson.D{{"create_time", -1}})

	filter := GetSearchFilter(searchParams)

	documentsCount, err := CollectionRecordRegister.CountDocuments(ctx, filter, nil)
	if err != nil {
		return nil, nil, err
	}
	meta := entity.Meta{
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    documentsCount,
	}

	cursor, err := CollectionRecordRegister.Find(ctx, filter, opt)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)
	records := make([]entity.SIPRecordRegister, 0, searchParams.PageSize)
	// 处理查询结果
	for cursor.Next(ctx) {
		var result entity.SIPRecordRegister
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		records = append(records, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, nil, err
	}
	return records, &meta, nil
}

func GetRecordCallList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error) {
	if searchParams.PageSize <= 0 {
		searchParams.PageSize = 10
	}
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}

	opt := options.Find().SetLimit(searchParams.PageSize).
		SetSkip(searchParams.PageSize * (searchParams.Page - 1)).
		SetSort(bson.D{{"create_time", -1}})

	filter := GetSearchFilter(searchParams)

	documentsCount, err := CollectionRecordCall.CountDocuments(ctx, filter, nil)
	if err != nil {
		return nil, nil, err
	}
	meta := entity.Meta{
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    documentsCount,
	}

	cursor, err := CollectionRecordCall.Find(ctx, filter, opt)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)
	records := make([]entity.SIPRecordCall, 0, searchParams.PageSize)
	// 处理查询结果
	for cursor.Next(ctx) {
		var result entity.SIPRecordCall
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		records = append(records, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, nil, err
	}
	return records, &meta, nil
}
