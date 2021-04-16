package handler

import (
	"net/http"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"gorm.io/gorm"
)

type UserHandler struct {
	service *service.UserService
}

func (handler *UserHandler) GetAllUsers(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)
	err, response := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

var userService *UserHandler
var onceUserService sync.Once

func GetUserHandler(db *gorm.DB) *UserHandler {
	onceUserService.Do(func() {
		userService = &UserHandler{service.GetUserService(db)}
	})

	return userService
}
