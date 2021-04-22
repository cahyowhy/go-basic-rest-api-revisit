package service

import (
	"sync"

	"gorm.io/gorm"
)

type baseService struct {
	db *gorm.DB
}

func (service *baseService) create(dest interface{}, omits ...string) error {
	tx := service.db
	if len(omits) > 0 {
		tx = tx.Omit(omits...)
	}

	return tx.Save(dest).Error
}

func (service *baseService) update(dest interface{}, id int, omits ...string) error {
	tx := service.db
	if len(omits) > 0 {
		tx = tx.Omit(omits...)
	}

	return tx.Updates(dest).Error
}

func (service *baseService) findAll(dest interface{}, offset int, limit int, filter map[string]interface{}, omits ...string) error {
	tx := service.db.Offset(offset).Limit(limit)
	if filter != nil {
		tx = tx.Where(filter)
	}

	if len(omits) > 0 {
		tx = tx.Omit(omits...)
	}

	return tx.Find(dest).Error
}

func (service *baseService) findWhere(dest interface{}, cond interface{}, omits ...string) error {
	tx := service.db
	if len(omits) > 0 {
		tx = tx.Omit(omits...)
	}

	return tx.First(dest, cond).Error
}

func (service *baseService) count(dest *int64, model interface{}, filter map[string]interface{}) error {
	tx := service.db.Model(model)
	if filter != nil {
		tx = tx.Where(filter)
	}

	return tx.Count(dest).Error
}

var service *baseService
var onceBaseService sync.Once

func getBaseService(db *gorm.DB) *baseService {
	onceBaseService.Do(func() {
		service = &baseService{db}
	})

	return service
}
