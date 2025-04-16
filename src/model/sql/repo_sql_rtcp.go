package sql

import (
	"context"
	"errors"
	"sip-monitor/src/entity"
	"time"

	"gorm.io/gorm"
)

func (r *GormRepository) CreateRtcpReportRaw(ctx context.Context, record *entity.RtcpReportRaw) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *GormRepository) CreateRtcpReportRaws(ctx context.Context, records []*entity.RtcpReportRaw) error {
	return r.db.WithContext(ctx).Create(records).Error
}

func (r *GormRepository) GetRtcpReportRawByID(ctx context.Context, id int64) (*entity.RtcpReportRaw, error) {
	var record entity.RtcpReportRaw
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *GormRepository) GetRtcpReportRawByBySIPCallID(ctx context.Context, sipCallID string) ([]*entity.RtcpReportRaw, error) {
	var records []*entity.RtcpReportRaw
	err := r.db.WithContext(ctx).Where("sip_call_id = ?", sipCallID).Find(&records).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return records, nil
}

func (r *GormRepository) DeleteRtcpReportRaw(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.RtcpReportRaw{}).Error
}

func (r *GormRepository) CreateRtcpReport(ctx context.Context, record *entity.RtcpReport) error {
	if record.CreateTime.IsZero() {
		record.CreateTime = time.Now()
	}
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *GormRepository) DeleteRtcpReport(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.RtcpReport{}).Error
}

func (r *GormRepository) GetRtcpReportByID(ctx context.Context, id int64) (*entity.RtcpReport, error) {
	var record entity.RtcpReport
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *GormRepository) GetRtcpReportBySIPCallID(ctx context.Context, sipCallID string) (*entity.RtcpReport, error) {
	var record entity.RtcpReport
	err := r.db.WithContext(ctx).Where("sip_call_id = ?", sipCallID).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}
