package service

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type UserService struct {
	db   *gorm.DB
	base *baseService
}

func (userService *UserService) FindAll(offset int, limit int, filter map[string]interface{}) (map[string]interface{}, error) {
	users := []model.User{}
	if err := userService.base.findAll(&users, offset, limit, filter, "Password"); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", users), nil
}

func (userService *UserService) Count(filter map[string]interface{}) (map[string]interface{}, error) {
	var total int64
	if err := userService.base.count(&total, &model.User{}, filter); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", total), nil
}

func (userService *UserService) SaveSession(username string) (userSession model.UserSession, err error) {
	// 3 days
	expired := time.Now().Local().Add(time.Hour * 24 * 3)
	refreshToken, err := util.GetUUID()

	if err != nil {
		return
	}

	if err = userService.db.First(&userSession, "refresh_token = ?", refreshToken).Error; err != nil {
		if err = userService.db.Where("username = ?", username).Delete(&model.UserSession{}).Error; err != nil {
			return
		}

		userSession.Expired = expired
		userSession.RefreshToken = refreshToken
		userSession.Username = username
		if err = userService.db.Save(&userSession).Error; err != nil {
			return
		}

		return userSession, nil
	}

	if err = userService.db.Model(&userSession).Update("expired", expired).Error; err != nil {
		return
	}

	return
}

func (userService *UserService) FindSession(refreshToken string) (map[string]interface{}, error) {
	userSession := model.UserSession{}
	if err := userService.db.First(&userSession, "refresh_token = ?", refreshToken).Error; err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	valid := userSession.Expired.Unix() > time.Now().Unix()
	if !valid {
		err := errors.New("user session not found")

		return util.ToMapKey("message", err.Error()), err
	}

	user := model.User{}
	if err := userService.db.First(&user, "username = ?", userSession.Username).Error; err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return responseWithToken(&user)
}

func (userService *UserService) Find(id int) (map[string]interface{}, error) {
	user := model.User{}

	if err := userService.base.findWhere(&user, id); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", user), nil
}

func (userService *UserService) Update(id int, body []byte) (map[string]interface{}, error) {
	user := model.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	user.Model.ID = uint(id)
	if err := userService.base.update(&user, id, "Password"); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", user), nil
}

func (userService *UserService) Create(body []byte) (map[string]interface{}, error) {
	user := model.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	if err := config.GetValidator().Struct(user); err != nil {
		return util.ToMapKey("errors", util.ValidationErrToString(err.(validator.ValidationErrors))), err
	}

	password, err := util.GeneratePassword(user.Password)
	if err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	user.Password = password
	if err := userService.base.create(&user); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", user), nil
}

func (userService *UserService) Logout(refreshToken string) error {
	return userService.db.Where("refresh_token = ?", refreshToken).Delete(&model.UserSession{}).Error
}

func (userService *UserService) Login(body []byte) (map[string]interface{}, string, error) {
	loginUser := model.LoginUser{}

	if err := json.Unmarshal(body, &loginUser); err != nil {
		return util.ToMapKey("message", err.Error()), "", err
	}

	if err := config.GetValidator().Struct(loginUser); err != nil {
		return util.ToMapKey("errors", util.ValidationErrToString(err.(validator.ValidationErrors))), "", err
	}

	var tx *gorm.DB
	user := model.User{}

	if len(loginUser.Email) > 0 {
		tx = userService.db.Where("email = ?", loginUser.Email).First(&user)
	} else if len(loginUser.Username) > 0 {
		tx = userService.db.Where("username = ?", loginUser.Username).First(&user)
	}

	if err := tx.Error; err != nil {
		return util.ToMapKey("message", err.Error()), "", err
	}

	isMatch := util.CompareHashPassword(loginUser.Password, user.Password)
	if !isMatch {
		err := errors.New("password are invalid")

		return util.ToMapKey("message", err.Error()), "", err
	}

	result, err := responseWithToken(&user)

	return result, user.Username, err
}

func responseWithToken(user *model.User) (map[string]interface{}, error) {
	token, err := util.GenerateJwt(*user)
	if err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	var response = make(map[string]interface{})
	userByte, err := json.Marshal(user)

	if err != nil {
		return util.ToMapKey("message", err.Error()), err
	} else if err := json.Unmarshal(userByte, &response); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	response["token"] = token
	return util.ToMapKey("data", response), nil
}

var userService *UserService
var onceUserService sync.Once

func GetUserService(db *gorm.DB) *UserService {
	onceUserService.Do(func() {
		userService = &UserService{db, getBaseService(db)}
	})

	return userService
}
