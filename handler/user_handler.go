package handler

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserHandler struct {
	service *service.UserService
}

func (handler *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)
	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserHandler) GetAllWithUserBook(w http.ResponseWriter, r *http.Request) {
	id, errParse := strconv.ParseInt(mux.Vars(r)["id"], 10, 8)
	if errParse != nil {
		util.ResponseSendJson(w, util.GetReponseMessage(errParse.Error()), http.StatusInternalServerError)

		return
	}

	response, err := handler.service.FindAllWithUserBook(int(id))

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (handler *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	response, err := handler.service.Create(r.Body)
	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
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

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	response, username, err := handler.service.Login(r.Body)
	if err != nil {
		util.ResponseSendJson(w, response, http.StatusUnauthorized)

		return
	} else if len(username) <= 0 {
		util.ResponseSendJson(w, util.GetReponseMessage("user not found"), http.StatusUnauthorized)

		return
	}

	handler.setRefreshCookie(w, username)
	util.ResponseSendJson(w, response)
}

func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		util.ResponseSendJson(w, util.GetReponseMessage(err.Error()), http.StatusUnauthorized)

		return
	} else if len(refreshToken.Value) <= 0 {
		util.ResponseSendJson(w, util.GetReponseMessage("refresh token not found"), http.StatusUnauthorized)

		return
	}

	if err := handler.service.Logout(refreshToken.Value); err != nil {
		util.ResponseSendJson(w, util.GetReponseMessage(err.Error()), http.StatusUnauthorized)

		return
	}

	cookie := http.Cookie{Name: "refresh_token", Value: "", Expires: time.Now(),
		HttpOnly: true, Secure: true, SameSite: http.SameSiteNoneMode}

	http.SetCookie(w, &cookie)
	util.ResponseSendJson(w, util.GetReponseMessage("logged out succeed"))
}

func (handler *UserHandler) Session(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		util.ResponseSendJson(w, util.GetReponseMessage(err.Error()), http.StatusUnauthorized)

		return
	}

	response, err := handler.service.FindSession(refreshToken.Value)
	if err != nil {
		util.ResponseSendJson(w, util.GetReponseMessage(err.Error()), http.StatusUnauthorized)

		return
	}

	util.ResponseSendJson(w, response)
}

func (handler *UserHandler) setRefreshCookie(w http.ResponseWriter, username string) {
	userSession, err := handler.service.SaveSession(username)
	if err != nil {
		util.ErrorLogger.Println(err)
		return
	}

	cookie := http.Cookie{Name: "refresh_token", Value: userSession.RefreshToken,
		Expires: userSession.Expired, HttpOnly: true, Secure: true, SameSite: http.SameSiteNoneMode}

	http.SetCookie(w, &cookie)
}

var userHandler *UserHandler
var onceUserHandler sync.Once

func GetUserHandler(db *gorm.DB) *UserHandler {
	onceUserHandler.Do(func() {
		userHandler = &UserHandler{service.GetUserService(db)}
	})

	return userHandler
}
