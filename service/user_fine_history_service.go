package service

import (
	"encoding/json"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"gorm.io/gorm"
)

type UserFineHistoryService struct {
	db   *gorm.DB
	base *baseService
}

func (service *UserFineHistoryService) FindAll(offset int, limit int, filter map[string]interface{}) (map[string]interface{}, error) {
	userFineHistory := []model.UserFineHistory{}

	tx := service.db.Omit("User.Password").Joins("User").Offset(offset).Limit(limit)
	if filter != nil {
		tx = tx.Where(filter)
	}

	if err := tx.Find(&userFineHistory).Error; err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", userFineHistory), nil
}

func (service *UserFineHistoryService) PayBookFine(userId int, body []byte) error {
	payFines := model.PayFines{}
	if err := json.Unmarshal(body, &payFines); err != nil {
		return err
	}

	ids := []uint{}
	for _, fine := range payFines.Fines {
		ids = append(ids, fine.ID)
	}

	filter := make(map[string]interface{})
	filter["user_id"] = userId
	filter["id"] = ids

	err := service.db.Model(model.UserFineHistory{}).Where(filter).Updates(map[string]interface{}{"has_paid": true}).Error
	if err != nil {
		return err
	}

	return nil
}

func (service *UserFineHistoryService) Count(filter map[string]interface{}) (map[string]interface{}, error) {
	var total int64
	if err := service.base.count(&total, &model.UserFineHistory{}, filter); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", total), nil
}

var userFineHistoryService *UserFineHistoryService
var onceUserFineHistoryService sync.Once

func GetUserFineHistoryService(db *gorm.DB) *UserFineHistoryService {
	onceUserFineHistoryService.Do(func() {
		userFineHistoryService = &UserFineHistoryService{db, getBaseService(db)}
	})

	return userFineHistoryService
}
