package sql

import (
	"context"
	"errors"
	"sip-monitor/src/entity"
	"time"

	"gorm.io/gorm"
)

func (r *GormRepository) GetCallIDsBySessionID(ctx context.Context, sessionID string) ([]string, error) {
	var sipCallIDs []string
	err := r.db.WithContext(ctx).Model(&entity.Call{}).Where("session_id = ?", sessionID).Distinct().Pluck("sip_call_id", &sipCallIDs).Error
	if err != nil {
		return nil, err
	}
	return sipCallIDs, nil
}

func (r *GormRepository) CreateCall(ctx context.Context, record *entity.Call) error {
	if record.CreateTime.IsZero() {
		now := time.Now()
		record.CreateTime = &now
	}
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *GormRepository) GetCallByID(ctx context.Context, id string) (*entity.Call, error) {
	var record entity.Call
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *GormRepository) GetCallBySIPCallID(ctx context.Context, sipCallID string) (*entity.Call, error) {
	var record entity.Call
	err := r.db.WithContext(ctx).Where("sip_call_id = ?", sipCallID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *GormRepository) GetCallList(ctx context.Context, params entity.SearchParams) ([]entity.Call, *entity.Meta, error) {
	var records []entity.Call
	var totalCount int64

	query := r.db.WithContext(ctx)

	if params.BeginTime != nil && params.EndTime != nil {
		query = query.Where("create_time BETWEEN ? AND ?", params.BeginTime, params.EndTime)
	}

	if params.SipCallID != "" {
		query = query.Where("sip_call_id = ?", params.SipCallID)
	}

	if params.SessionID != "" {
		query = query.Where("session_id = ?", params.SessionID)
	}

	if params.FromUser != "" {
		query = query.Where("from_user LIKE ?", "%"+params.FromUser+"%")
	}

	if params.ToUser != "" {
		query = query.Where("to_user LIKE ?", "%"+params.ToUser+"%")
	}

	if params.SrcHost != "" {
		query = query.Where("src_addr = ?", params.SrcHost)
	}

	if params.DstHost != "" {
		query = query.Where("dst_addr = ?", params.DstHost)
	}

	if params.HangupCode != "" {
		query = query.Where("hangup_code = ?", params.HangupCode)
	}

	// Count total records
	err := query.Model(&entity.Call{}).Count(&totalCount).Error
	if err != nil {
		return nil, nil, err
	}

	// Apply pagination
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(int(offset)).Limit(int(params.PageSize))
	}

	// Apply sorting
	if params.SortBy != "" {
		direction := "ASC"
		if params.SortDesc {
			direction = "DESC"
		}
		query = query.Order(params.SortBy + " " + direction)
	} else {
		// Default sort by creation time descending
		query = query.Order("create_time DESC")
	}

	// Execute query
	err = query.Find(&records).Error
	if err != nil {
		return nil, nil, err
	}

	// Calculate pagination metadata
	meta := r.calculatePagination(totalCount, int(params.Page), int(params.PageSize))

	return records, meta, nil
}

func (r *GormRepository) DeleteCall(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Call{}).Error
}
