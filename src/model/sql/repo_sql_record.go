package sql

import (
	"context"
	"errors"
	"sip-monitor/src/entity"
	"time"

	"gorm.io/gorm"
)

func (r *GormRepository) CreateRecordRaw(ctx context.Context, record *entity.RecordRaw) error {
	return r.db.WithContext(ctx).Create(record).Error
}

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

func (r *GormRepository) DeleteRecordRaw(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.RecordRaw{}).Error
}

func (r *GormRepository) CreateRecord(ctx context.Context, record *entity.Record) error {
	if record.CreateTime.IsZero() {
		record.CreateTime = time.Now()
	}
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *GormRepository) DeleteRecord(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Record{}).Error
}

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

func (r *GormRepository) GetRecordsBySIPCallIDs(ctx context.Context, sipCallIDs []string) ([]entity.Record, error) {
	var records []entity.Record
	var err error

	if len(sipCallIDs) == 1 {
		err = r.db.WithContext(ctx).Where("sip_call_id = ?", sipCallIDs[0]).Order("timestamp_micro").Find(&records).Error
	} else {
		err = r.db.WithContext(ctx).Where("sip_call_id IN ?", sipCallIDs).Order("timestamp_micro").Find(&records).Error
	}

	if err != nil {
		return nil, err
	}
	return records, nil
}
