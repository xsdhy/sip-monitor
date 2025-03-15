package model

import (
	"context"
	"errors"
	"sip-monitor/src/entity"
	"time"

	"gorm.io/gorm"
)

// GormRepository implements Repository using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new repository instance
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// buildSearchQuery constructs a query based on search parameters
func (r *GormRepository) buildSearchQuery(db *gorm.DB, params entity.SearchParams) *gorm.DB {
	query := db

	if params.BeginTime != nil && params.EndTime != nil {
		query = query.Where("create_time BETWEEN ? AND ?", params.BeginTime, params.EndTime)
	}

	if params.NodeIP != "" {
		query = query.Where("node_ip = ?", params.NodeIP)
	}

	if params.SipCallID != "" {
		query = query.Where("sip_call_id = ?", params.SipCallID)
	}

	if params.UserAgent != "" {
		query = query.Where("user_agent LIKE ?", "%"+params.UserAgent+"%")
	}

	if params.FromUser != "" {
		query = query.Where("from_user LIKE ?", "%"+params.FromUser+"%")
	}

	if params.SrcHost != "" {
		query = query.Where("src_host = ?", params.SrcHost)
	}

	if params.ToUser != "" {
		query = query.Where("to_user LIKE ?", "%"+params.ToUser+"%")
	}

	if params.DstHost != "" {
		query = query.Where("dst_host = ?", params.DstHost)
	}

	return query
}

// calculatePagination calculates pagination metrics
func (r *GormRepository) calculatePagination(totalCount int64, page, pageSize int) *entity.Meta {
	meta := &entity.Meta{
		Total:    totalCount,
		PageSize: int64(pageSize),
	}

	return meta
}

// Record operations

// CreateRecord creates a new record
func (r *GormRepository) CreateRecord(ctx context.Context, record *entity.Record) error {
	if record.CreateTime.IsZero() {
		record.CreateTime = time.Now()
	}
	return r.db.WithContext(ctx).Create(record).Error
}

// GetRecordByID retrieves a record by ID
func (r *GormRepository) GetRecordByID(ctx context.Context, id string) (*entity.Record, error) {
	var record entity.Record
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetRecordList retrieves records based on search parameters
func (r *GormRepository) GetRecordList(ctx context.Context, params entity.SearchParams) ([]entity.Record, *entity.Meta, error) {
	var records []entity.Record
	var totalCount int64

	// Build query
	query := r.buildSearchQuery(r.db.WithContext(ctx), params)

	// Count total records
	err := query.Model(&entity.Record{}).Count(&totalCount).Error
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
		query = query.Order("created_at DESC")
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

// DeleteRecord deletes a record
func (r *GormRepository) DeleteRecord(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Record{}).Error
}

func (r *GormRepository) GetRecordsBySIPCallID(ctx context.Context, sipCallID string) ([]entity.Record, error) {
	var records []entity.Record
	err := r.db.WithContext(ctx).Where("sip_call_id = ?", sipCallID).Order("timestamp_micro").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// SIP Call record operations

// CreateSIPCallRecord creates a new SIP call record
func (r *GormRepository) CreateSIPCallRecord(ctx context.Context, record *entity.SIPRecordCall) error {
	if record.CreateTime.IsZero() {
		now := time.Now()
		record.CreateTime = &now
	}
	return r.db.WithContext(ctx).Create(record).Error
}

// GetSIPCallRecordByID retrieves a SIP call record by ID
func (r *GormRepository) GetSIPCallRecordByID(ctx context.Context, id string) (*entity.SIPRecordCall, error) {
	var record entity.SIPRecordCall
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetSIPCallRecordList retrieves SIP call records based on search parameters
func (r *GormRepository) GetSIPCallRecordList(ctx context.Context, params entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error) {
	var records []entity.SIPRecordCall
	var totalCount int64

	// Build query
	query := r.buildSearchQuery(r.db.WithContext(ctx), params)

	// Count total records
	err := query.Model(&entity.SIPRecordCall{}).Count(&totalCount).Error
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

// DeleteSIPCallRecord deletes a SIP call record
func (r *GormRepository) DeleteSIPCallRecord(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.SIPRecordCall{}).Error
}

// SIP Register record operations

// CreateSIPRegisterRecord creates a new SIP register record
func (r *GormRepository) CreateSIPRegisterRecord(ctx context.Context, record *entity.SIPRecordRegister) error {
	if record.CreateTime.IsZero() {
		now := time.Now()
		record.CreateTime = &now
	}
	return r.db.WithContext(ctx).Create(record).Error
}

// GetSIPRegisterRecordByID retrieves a SIP register record by ID
func (r *GormRepository) GetSIPRegisterRecordByID(ctx context.Context, id string) (*entity.SIPRecordRegister, error) {
	var record entity.SIPRecordRegister
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetSIPRegisterRecordList retrieves SIP register records based on search parameters
func (r *GormRepository) GetSIPRegisterRecordList(ctx context.Context, params entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error) {
	var records []entity.SIPRecordRegister
	var totalCount int64

	// Build query
	query := r.buildSearchQuery(r.db.WithContext(ctx), params)

	// Count total records
	err := query.Model(&entity.SIPRecordRegister{}).Count(&totalCount).Error
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
		query = query.Order("created_at DESC")
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

// DeleteSIPRegisterRecord deletes a SIP register record
func (r *GormRepository) DeleteSIPRegisterRecord(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.SIPRecordRegister{}).Error
}
