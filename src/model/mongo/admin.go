package mongo

import (
	"context"
	"fmt"

	"sip-monitor/src/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MgInfra) GetSearchFilter(sp entity.SearchParams) bson.M {
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
	if sp.CallID != "" {
		filter["call_id"] = sp.CallID
	}

	if sp.UserAgent != "" {
		if sp.UserAgentOpr == "neq" {
			filter["user_agent"] = bson.M{"$regex": fmt.Sprintf(MgWhereNotContain, sp.UserAgent)}
		} else {
			filter["user_agent"] = bson.M{"$regex": fmt.Sprintf(MgWhereContain, sp.UserAgent)}
		}
	}

	if sp.FromUser != "" {
		if sp.FromUserOpr == "neq" {
			filter["from_user"] = bson.M{"$regex": fmt.Sprintf(MgWhereNotContain, sp.FromUser)}
		} else {
			filter["from_user"] = bson.M{"$regex": fmt.Sprintf(MgWhereContain, sp.FromUser)}
		}
	}

	if sp.SrcHost != "" {
		if sp.SrcHostOpr == "neq" {
			filter["scr_host"] = bson.M{"$ne": sp.SrcHost}
		} else {
			filter["scr_host"] = sp.SrcHost
		}
	}

	if sp.SrcCountryName != "" {
		if sp.SrcCountryNameOpr == "neq" {
			filter["src_country_name"] = bson.M{"$regex": fmt.Sprintf(MgWhereNotContain, sp.SrcCountryName)}
		} else {
			filter["src_country_name"] = bson.M{"$regex": fmt.Sprintf(MgWhereContain, sp.SrcCountryName)}
		}
	}

	if sp.SrcCityName != "" {
		if sp.SrcCityNameOpr == "neq" {
			filter["src_city_name"] = bson.M{"$regex": fmt.Sprintf(MgWhereNotContain, sp.SrcCityName)}
		} else {
			filter["src_city_name"] = bson.M{"$regex": fmt.Sprintf(MgWhereContain, sp.SrcCityName)}
		}
	}

	if sp.ToUser != "" {
		if sp.ToUserOpr == "neq" {
			filter["to_user"] = bson.M{"$regex": fmt.Sprintf(MgWhereNotContain, sp.ToUser)}
		} else {
			filter["to_user"] = bson.M{"$regex": fmt.Sprintf(MgWhereContain, sp.ToUser)}
		}
	}

	if sp.DstHost != "" {
		if sp.DstHostOpr == "neq" {
			filter["dst_host"] = bson.M{"$ne": sp.DstHost}
		} else {
			filter["dst_host"] = sp.DstHost
		}
	}

	if sp.DstCountryName != "" {
		if sp.DstCountryNameOpr == "neq" {
			filter["dst_country_name"] = bson.M{"$regex": fmt.Sprintf(MgWhereNotContain, sp.DstCountryName)}
		} else {
			filter["dst_country_name"] = bson.M{"$regex": fmt.Sprintf(MgWhereContain, sp.DstCountryName)}
		}
	}

	if sp.DstCityName != "" {
		if sp.DstCityNameOpr == "neq" {
			filter["dst_city_name"] = bson.M{"$regex": fmt.Sprintf(MgWhereNotContain, sp.DstCityName)}
		} else {
			filter["dst_city_name"] = bson.M{"$regex": fmt.Sprintf(MgWhereContain, sp.DstCityName)}
		}
	}

	return filter
}

func (m *MgInfra) GetDetailsBySipCallID(ctx context.Context, searchParams entity.SearchParams) ([]entity.Record, error) {
	opt := options.Find().SetSort(bson.D{
		{"create_time", 1},
		{"timestamp_micro", 1},
	})
	filter := m.GetSearchFilter(searchParams)

	cursor, err := m.CollectionRecord.Find(ctx, filter, opt)
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

func (m *MgInfra) GetRecordRegisterList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error) {
	if searchParams.PageSize <= 0 {
		searchParams.PageSize = 10
	}
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}

	opt := options.Find().SetLimit(searchParams.PageSize).
		SetSkip(searchParams.PageSize * (searchParams.Page - 1)).
		SetSort(bson.D{{"create_time", -1}})

	filter := m.GetSearchFilter(searchParams)

	documentsCount, err := m.CollectionRecordRegister.CountDocuments(ctx, filter, nil)
	if err != nil {
		return nil, nil, err
	}
	meta := entity.Meta{
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    documentsCount,
	}

	cursor, err := m.CollectionRecordRegister.Find(ctx, filter, opt)
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

func (m *MgInfra) GetRecordCallList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error) {
	if searchParams.PageSize <= 0 {
		searchParams.PageSize = 10
	}
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}

	opt := options.Find().SetLimit(searchParams.PageSize).
		SetSkip(searchParams.PageSize * (searchParams.Page - 1)).
		SetSort(bson.D{{"create_time", -1}})

	filter := m.GetSearchFilter(searchParams)

	documentsCount, err := m.CollectionRecordCall.CountDocuments(ctx, filter, nil)
	if err != nil {
		return nil, nil, err
	}
	meta := entity.Meta{
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    documentsCount,
	}

	cursor, err := m.CollectionRecordCall.Find(ctx, filter, opt)
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
