package sql

import (
	"sip-monitor/src/entity"

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

// calculatePagination calculates pagination metrics
func (r *GormRepository) calculatePagination(totalCount int64, page, pageSize int) *entity.Meta {
	meta := &entity.Meta{
		Total:    totalCount,
		PageSize: int64(pageSize),
	}

	return meta
}
