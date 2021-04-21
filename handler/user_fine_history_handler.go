package handler

import (
	"net/http"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserFineHistoryHandler struct {
	service *service.UserFineHistoryService
}

func (handler *UserFineHistoryHandler) PayBookFine(w http.ResponseWriter, r *http.Request) {
	id, ok := util.ToInt(mux.Vars(r)["id"])
	if !ok {
		util.ResponseSendJson(w, util.ToMapKey("message", "invalid path params"), http.StatusInternalServerError)

		return
	}

	if err := handler.service.PayBookFine(id, r.Body); err != nil {
		util.ResponseSendJson(w, util.ToMapKey("message", err.Error()), http.StatusInternalServerError)

		return
	}

	util.ResponseSendJson(w, util.ToMapKey("message", "Success paid fine"))
}

func (handler *UserFineHistoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)
	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserFineHistoryHandler) Count(w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)
	response, err := handler.service.Count(queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

var userFineHistoryHandler *UserFineHistoryHandler
var onceUserFineHistoryHandler sync.Once

func GetUserFineHistoryHandler(db *gorm.DB) *UserFineHistoryHandler {
	onceUserFineHistoryHandler.Do(func() {
		userFineHistoryHandler = &UserFineHistoryHandler{service.GetUserFineHistoryService(db)}
	})

	return userFineHistoryHandler
}
