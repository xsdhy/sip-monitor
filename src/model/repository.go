package model

import (
	"context"
	"sip-monitor/src/entity"
)

// Repository defines the interface for database operations
type Repository interface {
	// Record raw operations
	CreateRecordRaw(ctx context.Context, record *entity.RecordRaw) error
	GetRecordRawByID(ctx context.Context, id int64) (*entity.RecordRaw, error)
	DeleteRecordRaw(ctx context.Context, id int64) error

	// Record operations
	CreateRecord(ctx context.Context, record *entity.Record) error
	DeleteRecord(ctx context.Context, id int64) error
	GetRecordsBySIPCallID(ctx context.Context, sipCallID string) ([]entity.Record, error)
	GetRecordsBySIPCallIDs(ctx context.Context, sipCallIDs []string) ([]entity.Record, error)

	// SIP Call record operations
	CreateSIPCallRecord(ctx context.Context, record *entity.SIPRecordCall) error
	GetSIPCallRecordByID(ctx context.Context, id string) (*entity.SIPRecordCall, error)
	GetSIPCallRecordBySIPCallID(ctx context.Context, sipCallID string) (*entity.SIPRecordCall, error)
	GetSIPCallIDsBySessionID(ctx context.Context, sessionID string) ([]string, error)
	GetSIPCallRecordList(ctx context.Context, params entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error)
	DeleteSIPCallRecord(ctx context.Context, id string) error

	// SIP Register record operations
	CreateSIPRegisterRecord(ctx context.Context, record *entity.SIPRecordRegister) error
	GetSIPRegisterRecordByID(ctx context.Context, id string) (*entity.SIPRecordRegister, error)
	GetSIPRegisterRecordList(ctx context.Context, params entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error)
	DeleteSIPRegisterRecord(ctx context.Context, id string) error

	// User operations
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetUsers(ctx context.Context) ([]entity.User, error)

	// Create default admin user
	CreateDefaultAdminUser(ctx context.Context) error
}
