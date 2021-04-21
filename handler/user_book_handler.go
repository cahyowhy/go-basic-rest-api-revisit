package handler

import (
	"net/http"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserBookHandler struct {
	service *service.UserBookService
}

func (handler *UserBookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)
	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserBookHandler) GetAllFromAuth(w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)

	claims, okClaim := r.Context().Value(util.KeyUser).(jwt.MapClaims)
	if !okClaim {
		util.ResponseSendJson(w, util.ToMapKey("message", "Unauthorize"), http.StatusUnauthorized)

		return
	}

	var id int
	idClaim, ok := claims["ID"]

	if ok {
		id, ok = util.ToInt(idClaim)
	}

	if !ok {
		util.ResponseSendJson(w, "invalid user id", http.StatusInternalServerError)

		return
	}

	if queryParam.Filter == nil {
		queryParam.Filter = make(map[string]interface{})
		queryParam.Filter["user_id"] = uint(id)
	}

	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserBookHandler) BorrowBooks(w http.ResponseWriter, r *http.Request) {
	id, ok := util.ToInt(mux.Vars(r)["id"])
	if !ok {
		util.ResponseSendJson(w, util.ToMapKey("message", "invalid path params"), http.StatusInternalServerError)

		return
	}

	response, err := handler.service.BorrowBook(uint(id), r.Body)
	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserBookHandler) ReturnBooks(w http.ResponseWriter, r *http.Request) {
	id, ok := util.ToInt(mux.Vars(r)["id"])
	if !ok {
		util.ResponseSendJson(w, util.ToMapKey("message", "invalid path params"), http.StatusInternalServerError)

		return
	}

	response, err := handler.service.ReturnBook(uint(id), r.Body)
	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserBookHandler) Count(w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)
	response, err := handler.service.Count(queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

var userBookHandler *UserBookHandler
var onceUserBookHandler sync.Once

func GetUserBookHandler(db *gorm.DB) *UserBookHandler {
	onceUserBookHandler.Do(func() {
		userBookHandler = &UserBookHandler{service.GetUserBookService(db)}
	})

	return userBookHandler
}
