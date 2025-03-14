package model

import (
	"context"
	"sip-monitor/src/entity"
)

// Repository defines the interface for database operations
type Repository interface {
	// Record operations
	CreateRecord(ctx context.Context, record *entity.Record) error
	GetRecordByID(ctx context.Context, id string) (*entity.Record, error)
	GetRecordList(ctx context.Context, params entity.SearchParams) ([]entity.Record, *entity.Meta, error)
	DeleteRecord(ctx context.Context, id string) error
	GetRecordsBySIPCallID(ctx context.Context, sipCallID string) ([]entity.Record, error)

	// SIP Call record operations
	CreateSIPCallRecord(ctx context.Context, record *entity.SIPRecordCall) error
	GetSIPCallRecordByID(ctx context.Context, id string) (*entity.SIPRecordCall, error)
	GetSIPCallRecordList(ctx context.Context, params entity.SearchParams) ([]entity.SIPRecordCall, *entity.Meta, error)
	DeleteSIPCallRecord(ctx context.Context, id string) error

	// SIP Register record operations
	CreateSIPRegisterRecord(ctx context.Context, record *entity.SIPRecordRegister) error
	GetSIPRegisterRecordByID(ctx context.Context, id string) (*entity.SIPRecordRegister, error)
	GetSIPRegisterRecordList(ctx context.Context, params entity.SearchParams) ([]entity.SIPRecordRegister, *entity.Meta, error)
	DeleteSIPRegisterRecord(ctx context.Context, id string) error
}
