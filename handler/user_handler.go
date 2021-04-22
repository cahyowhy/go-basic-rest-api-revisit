package handler

import (
	"net/http"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	service *service.UserService
}

func (handler *UserHandler) GetAll(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)
	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *UserHandler) Count(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)
	response, err := handler.service.Count(queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *UserHandler) Get(c *fiber.Ctx) error {
	id, ok := util.ToInt(c.Params("id"))

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(util.ToMapKey("message", "invalid path params"))
	}

	response, err := handler.service.Find(int(id))
	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusNotFound).JSON(response)
}

func (handler *UserHandler) Create(c *fiber.Ctx) error {
	response, err := handler.service.Create(c.Body())
	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *UserHandler) Update(c *fiber.Ctx) error {
	id, ok := util.ToInt(c.Params("id"))
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(util.ToMapKey("message", "invalid path params"))
	}

	response, err := handler.service.Update(int(id), c.Body())
	if err == nil {
		return c.Status(http.StatusInternalServerError).JSON(util.ToMapKey("message", "invalid path params"))
	}

	return c.JSON(response)
}

func (handler *UserHandler) Login(c *fiber.Ctx) error {
	response, username, err := handler.service.Login(c.Body())
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(response)
	} else if len(username) <= 0 {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", "user not found"))
	}

	handler.setRefreshCookie(c, username)
	return c.JSON(response)
}

func (handler *UserHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if len(refreshToken) <= 0 {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", "refresh token not found"))
	}

	if err := handler.service.Logout(refreshToken); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", err.Error()))
	}

	c.ClearCookie("refresh_token")
	return c.JSON(util.ToMapKey("message", "logged out succeed"))
}

func (handler *UserHandler) Session(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if len(refreshToken) <= 0 {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", "refresh token not found"))
	}

	response, err := handler.service.FindSession(refreshToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", err.Error()))
	}

	return c.JSON(response)
}

func (handler *UserHandler) setRefreshCookie(c *fiber.Ctx, username string) {
	userSession, err := handler.service.SaveSession(username)
	if err != nil {
		util.ErrorLogger.Println(err)
		return
	}

	cookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    userSession.RefreshToken,
		Expires:  userSession.Expired,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "none",
	}

	c.Cookie(cookie)
}

var userHandler *UserHandler
var onceUserHandler sync.Once

func GetUserHandler(db *gorm.DB) *UserHandler {
	onceUserHandler.Do(func() {
		userHandler = &UserHandler{service.GetUserService(db)}
	})

	return userHandler
}
