package service

import (
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func (userService *UserService) FindAll(offset int, limit int, filter map[string]interface{}) (err error, response map[string]interface{}) {
	users := []model.User{}
	tx := userService.db.Offset(offset).Limit(limit)

	if filter != nil {
		tx = tx.Where(filter)
	}

	if err := tx.Find(&users).Error; err != nil {
		return err, util.GetReponseMessage(err.Error())
	}

	var body = make(map[string]interface{})
	body["data"] = users

	return nil, body
}

var userService *UserService
var onceUserService sync.Once

func GetUserService(db *gorm.DB) *UserService {
	onceUserService.Do(func() {
		userService = &UserService{db}
	})

	return userService
}
