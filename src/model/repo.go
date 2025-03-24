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
	GetRecordsBySIPCallIDs(ctx context.Context, sipCallIDs []string) ([]entity.Record, error)

	// SIP Call record operations
	CreateCall(ctx context.Context, record *entity.Call) error
	GetCallByID(ctx context.Context, id string) (*entity.Call, error)
	GetCallBySIPCallID(ctx context.Context, sipCallID string) (*entity.Call, error)
	GetCallIDsBySessionID(ctx context.Context, sessionID string) ([]string, error)
	GetCallList(ctx context.Context, params entity.SearchParams) ([]entity.Call, *entity.Meta, error)
	DeleteCall(ctx context.Context, id string) error

	// User operations
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetUsers(ctx context.Context) ([]entity.User, error)
	CreateDefaultAdminUser(ctx context.Context) error

	GetCallStat(ctx context.Context, params entity.CallStatDTO) ([]entity.CallStatVO, error)
}
