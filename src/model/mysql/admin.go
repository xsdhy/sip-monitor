package mysql

import (
	"context"
	"fmt"
	"sip-monitor/src/entity"
)

// GetSearchFilter 生成搜索过滤条件
func (n *NoSqlInfra) GetSearchFilter(sp entity.SearchParams) map[string]interface{} {
	filter := map[string]interface{}{}

	if sp.BeginTime != nil && sp.EndTime != nil {
		filter["create_time >= ?"] = sp.BeginTime
		filter["create_time <= ?"] = sp.EndTime
	}

	if sp.NodeIP != "" {
		filter["node_ip"] = sp.NodeIP
	}
	if sp.CallID != "" {
		filter["call_id"] = sp.CallID
	}

	if sp.UserAgent != "" {
		if sp.UserAgentOpr == "neq" {
			filter["user_agent NOT LIKE ?"] = fmt.Sprintf("%%%s%%", sp.UserAgent)
		} else {
			filter["user_agent LIKE ?"] = fmt.Sprintf("%%%s%%", sp.UserAgent)
		}
	}

	if sp.FromUser != "" {
		if sp.FromUserOpr == "neq" {
			filter["from_user NOT LIKE ?"] = fmt.Sprintf("%%%s%%", sp.FromUser)
		} else {
			filter["from_user LIKE ?"] = fmt.Sprintf("%%%s%%", sp.FromUser)
		}
	}

	if sp.SrcHost != "" {
		if sp.SrcHostOpr == "neq" {
			filter["src_host <> ?"] = sp.SrcHost
		} else {
			filter["src_host"] = sp.SrcHost
		}
	}

	if sp.SrcCountryName != "" {
		if sp.SrcCountryNameOpr == "neq" {
			filter["src_country_name NOT LIKE ?"] = fmt.Sprintf("%%%s%%", sp.SrcCountryName)
		} else {
			filter["src_country_name LIKE ?"] = fmt.Sprintf("%%%s%%", sp.SrcCountryName)
		}
	}

	if sp.SrcCityName != "" {
		if sp.SrcCityNameOpr == "neq" {
			filter["src_city_name NOT LIKE ?"] = fmt.Sprintf("%%%s%%", sp.SrcCityName)
		} else {
			filter["src_city_name LIKE ?"] = fmt.Sprintf("%%%s%%", sp.SrcCityName)
		}
	}

	if sp.ToUser != "" {
		if sp.ToUserOpr == "neq" {
			filter["to_user NOT LIKE ?"] = fmt.Sprintf("%%%s%%", sp.ToUser)
		} else {
			filter["to_user LIKE ?"] = fmt.Sprintf("%%%s%%", sp.ToUser)
		}
	}

	if sp.DstHost != "" {
		if sp.DstHostOpr == "neq" {
			filter["dst_host <> ?"] = sp.DstHost
		} else {
			filter["dst_host"] = sp.DstHost
		}
	}

	if sp.DstCountryName != "" {
		if sp.DstCountryNameOpr == "neq" {
			filter["dst_country_name NOT LIKE ?"] = fmt.Sprintf("%%%s%%", sp.DstCountryName)
		} else {
			filter["dst_country_name LIKE ?"] = fmt.Sprintf("%%%s%%", sp.DstCountryName)
		}
	}

	if sp.DstCityName != "" {
		if sp.DstCityNameOpr == "neq" {
			filter["dst_city_name NOT LIKE ?"] = fmt.Sprintf("%%%s%%", sp.DstCityName)
		} else {
			filter["dst_city_name LIKE ?"] = fmt.Sprintf("%%%s%%", sp.DstCityName)
		}
	}

	return filter
}

// GetDetailsBySipCallID 根据 SIP Call ID 获取详细信息
func (n *NoSqlInfra) GetDetailsBySipCallID(ctx context.Context, searchParams entity.SearchParams) ([]entity.Record, error) {
	var records []entity.Record
	filter := n.GetSearchFilter(searchParams)

	err := n.db.WithContext(ctx).
		Where(filter).
		Order("create_time ASC").
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	return records, nil
}

// GetRecordRegisterList 获取注册记录列表
func (n *NoSqlInfra) GetRecordRegisterList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error) {
	if searchParams.PageSize <= 0 {
		searchParams.PageSize = 10
	}
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}

	var records []entity.SIPRecordRegister
	filter := n.GetSearchFilter(searchParams)

	var total int64
	n.db.WithContext(ctx).Model(&entity.SIPRecordRegister{}).Where(filter).Count(&total)

	meta := &entity.Meta{
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    total,
	}

	err := n.db.WithContext(ctx).
		Where(filter).
		Order("create_time DESC").
		Limit(int(searchParams.PageSize)).
		Offset(int(searchParams.PageSize * (searchParams.Page - 1))).
		Find(&records).Error

	if err != nil {
		return nil, nil, err
	}

	return records, meta, nil
}

// GetRecordCallList 获取通话记录列表
func (n *NoSqlInfra) GetRecordCallList(ctx context.Context, searchParams entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error) {
	if searchParams.PageSize <= 0 {
		searchParams.PageSize = 10
	}
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}

	var records []entity.SIPRecordCall
	filter := n.GetSearchFilter(searchParams)

	var total int64
	n.db.WithContext(ctx).Model(&entity.SIPRecordCall{}).Where(filter).Count(&total)

	meta := &entity.Meta{
		Page:     searchParams.Page,
		PageSize: searchParams.PageSize,
		Total:    total,
	}

	err := n.db.WithContext(ctx).
		Where(filter).
		Order("create_time DESC").
		Limit(int(searchParams.PageSize)).
		Offset(int(searchParams.PageSize * (searchParams.Page - 1))).
		Find(&records).Error

	if err != nil {
		return nil, nil, err
	}

	return records, meta, nil
}
