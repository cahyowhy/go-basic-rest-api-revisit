package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/dgrijalva/jwt-go"
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

func (handler *UserHandler) Count(w http.ResponseWriter, r *http.Request) {
	queryParam := GetQueryParam(r)
	response, err := handler.service.Count(queryParam.Filter)

	if err == nil {
		util.ResponseSendJson(w, response)

		return
	}

	util.ResponseSendJson(w, response, http.StatusInternalServerError)
}

func (handler *UserHandler) GetAllWithUserBook(w http.ResponseWriter, r *http.Request) {
	id, ok := util.ToInt(mux.Vars(r)["id"])
	if !ok {
		util.ResponseSendJson(w, util.ToMapKey("message", "invalid path params"), http.StatusInternalServerError)

		return
	}

	handler.findUserBookByUserId(&w, int(id))
}

func (handler *UserHandler) GetAllWithUserBookFromToken(w http.ResponseWriter, r *http.Request) {
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

	handler.findUserBookByUserId(&w, id)
}

func (handler *UserHandler) findUserBookByUserId(w *http.ResponseWriter, id int) {
	response, err := handler.service.FindAllWithUserBook(id)

	if err == nil {
		util.ResponseSendJson(*w, response)

		return
	}

	util.ResponseSendJson(*w, response, http.StatusInternalServerError)
}

func (handler *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := util.ToInt(mux.Vars(r)["id"])
	if !ok {
		util.ResponseSendJson(w, util.ToMapKey("message", "invalid path params"), http.StatusInternalServerError)

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
	id, ok := util.ToInt(mux.Vars(r)["id"])
	if !ok {
		util.ResponseSendJson(w, util.ToMapKey("message", "invalid path params"), http.StatusInternalServerError)

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
		util.ResponseSendJson(w, util.ToMapKey("message", "user not found"), http.StatusUnauthorized)

		return
	}

	handler.setRefreshCookie(w, username)
	util.ResponseSendJson(w, response)
}

func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		util.ResponseSendJson(w, util.ToMapKey("message", err.Error()), http.StatusUnauthorized)

		return
	} else if len(refreshToken.Value) <= 0 {
		util.ResponseSendJson(w, util.ToMapKey("message", "refresh token not found"), http.StatusUnauthorized)

		return
	}

	if err := handler.service.Logout(refreshToken.Value); err != nil {
		util.ResponseSendJson(w, util.ToMapKey("message", err.Error()), http.StatusUnauthorized)

		return
	}

	cookie := http.Cookie{Name: "refresh_token", Value: "", Expires: time.Now(),
		HttpOnly: true, Secure: true, SameSite: http.SameSiteNoneMode}

	http.SetCookie(w, &cookie)
	util.ResponseSendJson(w, util.ToMapKey("message", "logged out succeed"))
}

func (handler *UserHandler) Session(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		util.ResponseSendJson(w, util.ToMapKey("message", err.Error()), http.StatusUnauthorized)

		return
	}

	response, err := handler.service.FindSession(refreshToken.Value)
	if err != nil {
		util.ResponseSendJson(w, util.ToMapKey("message", err.Error()), http.StatusUnauthorized)

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
