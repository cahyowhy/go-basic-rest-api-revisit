package handler

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type BookHandler struct {
	service *service.BookService
}

func (handler *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)
	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, errParse := strconv.ParseInt(mux.Vars(r)["id"], 10, 8)
	if errParse != nil {
		util.ResponseSendJson(w, util.GetReponseMessage(errParse.Error()), http.StatusInternalServerError)

		return
	}

	response, err := handler.service.Find(int(id))
	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusNotFound)
}

func (handler *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	response, err := handler.service.Create(r.Body)
	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, errParse := strconv.ParseInt(mux.Vars(r)["id"], 10, 8)
	if errParse != nil {
		util.ResponseSendJson(w, util.GetReponseMessage(errParse.Error()), http.StatusInternalServerError)

		return
	}

	response, err := handler.service.Update(int(id), r.Body)
	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

var bookHandler *BookHandler
var onceBookHandler sync.Once

func GetBookHandler(db *gorm.DB) *BookHandler {
	onceBookHandler.Do(func() {
		bookHandler = &BookHandler{service.GetBookService(db)}
	})

	return bookHandler
}
