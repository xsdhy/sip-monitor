package sql

import (
	"context"
	"sip-monitor/src/entity"
)

func (r *GormRepository) GetCallStat(ctx context.Context, params entity.CallStatDTO) ([]*entity.CallStatVO, error) {
	query := r.db.WithContext(ctx).Model(&entity.Call{}).
		Select("dst_addr as ip, COUNT(*) AS total, SUM(CASE WHEN talk_duration > 0 THEN 1 ELSE 0 END) AS answered, SUM(CASE WHEN hangup_code = 0 THEN 1 ELSE 0 END) AS hangup_code_0_count, SUM(CASE WHEN hangup_code BETWEEN 100 AND 199 THEN 1 ELSE 0 END) AS hangup_code_1xx_count, SUM(CASE WHEN hangup_code BETWEEN 200 AND 299 THEN 1 ELSE 0 END) AS hangup_code_2xx_count, SUM(CASE WHEN hangup_code BETWEEN 300 AND 399 THEN 1 ELSE 0 END) AS hangup_code_3xx_count, SUM(CASE WHEN hangup_code BETWEEN 400 AND 499 THEN 1 ELSE 0 END) AS hangup_code_4xx_count, SUM(CASE WHEN hangup_code BETWEEN 500 AND 599 THEN 1 ELSE 0 END) AS hangup_code_5xx_count").
		Group("dst_addr")

	if params.BeginTime != nil {
		query = query.Where("create_time >= ?", params.BeginTime)
	}

	if params.EndTime != nil {
		query = query.Where("create_time <= ?", params.EndTime)
	}

	var result []*entity.CallStatVO
	err := query.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
