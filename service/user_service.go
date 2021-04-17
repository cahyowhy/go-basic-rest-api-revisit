package service

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func (userService *UserService) FindAll(offset int, limit int, filter map[string]interface{}) (map[string]interface{}, error) {
	users := []model.User{}
	tx := userService.db.Offset(offset).Limit(limit)

	if filter != nil {
		tx = tx.Where(filter)
	}

	if err := tx.Omit("password").Find(&users).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var body = make(map[string]interface{})
	body["data"] = users

	return body, nil
}

func (userService *UserService) SaveSession(username string) (model.UserSession, error) {
	userSession := model.UserSession{}

	// 3 days
	expired := time.Now().Local().Add(time.Hour * 24 * 3)
	refreshToken, err := util.GetUUID()

	if err != nil {
		return userSession, err
	}

	if err := userService.db.First(&userSession, "refresh_token = ?", refreshToken).Error; err != nil {
		userSession.Expired = expired
		userSession.RefreshToken = refreshToken
		userSession.Username = username

		if err := userService.db.Save(&userSession).Error; err != nil {
			return userSession, err
		}

		return userSession, nil
	}

	if err := userService.db.Model(&userSession).Update("expired", expired).Error; err != nil {
		return userSession, err
	}

	return userSession, nil
}

func (userService *UserService) FindSession(refreshToken string) (map[string]interface{}, error) {
	userSession := model.UserSession{}

	if err := userService.db.First(&userSession, "refresh_token = ?", refreshToken).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	valid := userSession.Expired.Unix() > time.Now().Unix()
	if !valid {
		err := errors.New("user session not found")

		return util.GetReponseMessage(err.Error()), err
	}

	user := model.User{}
	if err := userService.db.First(&user, "username = ?", userSession.Username).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	return responseWithToken(&user)
}

func (userService *UserService) Find(id int) (map[string]interface{}, error) {
	user := model.User{}

	if err := userService.db.Omit("password").First(&user, id).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var body = make(map[string]interface{})
	body["data"] = user

	return body, nil
}

func (userService *UserService) Update(id int, body io.Reader) (map[string]interface{}, error) {
	user := model.User{}
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&user); err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	user.Model.ID = uint(id)

	if err := userService.db.Omit("password").Updates(&user).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var response = make(map[string]interface{})
	response["data"] = user

	return response, nil
}

func (userService *UserService) Create(body io.Reader) (map[string]interface{}, error) {
	user := model.User{}
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&user); err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var err = config.GetValidator().Struct(user)
	if err != nil {
		response := make(map[string]interface{})
		response["errors"] = util.ValidationErrToString(err.(validator.ValidationErrors))

		return response, err
	}

	password, err := util.GeneratePassword(user.Password)

	if err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	user.Password = password

	if err := userService.db.Save(&user).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var response = make(map[string]interface{})
	response["data"] = user

	return response, nil
}

func (userService *UserService) Logout(refreshToken string) error {
	return userService.db.Where("refresh_token = ?", refreshToken).Delete(&model.UserSession{}).Error
}

func (userService *UserService) Login(body io.Reader) (map[string]interface{}, string, error) {
	loginUser := model.LoginUser{}
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&loginUser); err != nil {
		return util.GetReponseMessage(err.Error()), "", err
	}

	var err = config.GetValidator().Struct(loginUser)
	if err != nil {
		response := make(map[string]interface{})
		response["errors"] = util.ValidationErrToString(err.(validator.ValidationErrors))

		return response, "", err
	}

	var tx *gorm.DB
	user := model.User{}

	if len(loginUser.Email) > 0 {
		tx = userService.db.Where("email = ?", loginUser.Email).First(&user)
	} else if len(loginUser.Username) > 0 {
		tx = userService.db.Where("username = ?", loginUser.Username).First(&user)
	}

	if err := tx.Error; err != nil {
		return util.GetReponseMessage(err.Error()), "", err
	}

	isMatch := util.CompareHashPassword(loginUser.Password, user.Password)
	if !isMatch {
		err := errors.New("password are invalid")

		return util.GetReponseMessage(err.Error()), "", err
	}

	result, err := responseWithToken(&user)

	return result, user.Username, err
}

func responseWithToken(user *model.User) (map[string]interface{}, error) {
	token, err := util.GenerateJwt(*user)
	if err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var response = make(map[string]interface{})
	userByte, err := json.Marshal(user)

	if err != nil {
		return util.GetReponseMessage(err.Error()), err
	} else if err := json.Unmarshal(userByte, &response); err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	response["token"] = token

	responseData := make(map[string]interface{})
	responseData["data"] = response

	return responseData, nil
}

var userService *UserService
var onceUserService sync.Once

func GetUserService(db *gorm.DB) *UserService {
	onceUserService.Do(func() {
		userService = &UserService{db}
	})

	return userService
}
