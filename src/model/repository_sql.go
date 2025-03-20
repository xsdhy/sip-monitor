package model

import (
	"context"
	"errors"
	"sip-monitor/src/entity"
	"time"

	"golang.org/x/crypto/bcrypt"
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

// Record raw operations

// CreateRecordRaw creates a new record raw
func (r *GormRepository) CreateRecordRaw(ctx context.Context, record *entity.RecordRaw) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// GetRecordRawByID retrieves a record raw by ID
func (r *GormRepository) GetRecordRawByID(ctx context.Context, id int64) (*entity.RecordRaw, error) {
	var record entity.RecordRaw
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// DeleteRecordRaw deletes a record raw
func (r *GormRepository) DeleteRecordRaw(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.RecordRaw{}).Error
}

// Record operations

// CreateRecord creates a new record
func (r *GormRepository) CreateRecord(ctx context.Context, record *entity.Record) error {
	if record.CreateTime.IsZero() {
		record.CreateTime = time.Now()
	}
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *GormRepository) DeleteRecord(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Record{}).Error
}

// GetRecordByID retrieves a record by ID
func (r *GormRepository) GetRecordByID(ctx context.Context, id int64) (*entity.Record, error) {
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

func (r *GormRepository) GetRecordsBySIPCallID(ctx context.Context, sipCallID string) ([]entity.Record, error) {
	var records []entity.Record
	err := r.db.WithContext(ctx).Where("sip_call_id = ?", sipCallID).Order("timestamp_micro").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *GormRepository) GetRecordsBySIPCallIDs(ctx context.Context, sipCallIDs []string) ([]entity.Record, error) {
	var records []entity.Record
	err := r.db.WithContext(ctx).Where("sip_call_id IN ?", sipCallIDs).Order("timestamp_micro").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *GormRepository) GetSIPCallIDsBySessionID(ctx context.Context, sessionID string) ([]string, error) {
	var sipCallIDs []string
	err := r.db.WithContext(ctx).Model(&entity.SIPRecordCall{}).Where("session_id = ?", sessionID).Distinct().Pluck("sip_call_id", &sipCallIDs).Error
	if err != nil {
		return nil, err
	}
	return sipCallIDs, nil
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

func (r *GormRepository) GetSIPCallRecordBySIPCallID(ctx context.Context, sipCallID string) (*entity.SIPRecordCall, error) {
	var record entity.SIPRecordCall
	err := r.db.WithContext(ctx).Where("sip_call_id = ?", sipCallID).First(&record).Error
	if err != nil {
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

// CreateUser creates a new user in the database
func (r *GormRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if user.CreateAt.IsZero() {
		user.CreateAt = time.Now()
	}
	user.UpdateAt = time.Now()
	return r.db.WithContext(ctx).Create(user).Error
}

// GetUserByID retrieves a user by ID
func (r *GormRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *GormRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user
func (r *GormRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdateAt = time.Now()
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(user).Error
}

// DeleteUser deletes a user by ID
func (r *GormRepository) DeleteUser(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *GormRepository) CreateDefaultAdminUser(ctx context.Context) error {
	var userNum int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&userNum).Error
	if err != nil {
		return err
	}

	if userNum > 0 {
		return nil
	}

	// Create default admin user

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create the admin user
	adminUser := &entity.User{
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "Administrator",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	// Save the admin user to the database
	if err := r.CreateUser(ctx, adminUser); err != nil {
		return err
	}

	return nil
}

func (r *GormRepository) GetUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}
