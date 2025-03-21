package sql

import (
	"sip-monitor/src/entity"
)

func (r *GormRepository) GatewayCreate(gateway *entity.Gateway) error {
	return r.db.Create(gateway).Error
}

func (r *GormRepository) GatewayGetByID(id int64) (*entity.Gateway, error) {
	var gateway entity.Gateway
	err := r.db.Where("id = ?", id).First(&gateway).Error
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}

func (r *GormRepository) GatewayList() ([]entity.Gateway, error) {
	var gateways []entity.Gateway
	err := r.db.Find(&gateways).Error
	if err != nil {
		return nil, err
	}
	return gateways, nil
}

func (r *GormRepository) GatewayUpdate(gateway *entity.Gateway) error {
	return r.db.Select("Name", "Addr", "Remark", "UpdateAt").Save(gateway).Error
}

func (r *GormRepository) GatewayDelete(id int64) error {
	return r.db.Delete(&entity.Gateway{}, id).Error
}

func (r *GormRepository) GatewayGetByName(name string) (*entity.Gateway, error) {
	var gateway entity.Gateway
	err := r.db.Where("name = ?", name).First(&gateway).Error
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}

func (r *GormRepository) GatewayGetByAddr(addr string) (*entity.Gateway, error) {
	var gateway entity.Gateway
	err := r.db.Where("addr = ?", addr).First(&gateway).Error
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}
