package handler

import (
	"net/http"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserBookHandler struct {
	service *service.UserBookService
}

func (handler *UserBookHandler) GetAll(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)
	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *UserBookHandler) Count(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)
	response, err := handler.service.Count(queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *UserBookHandler) GetAllFromAuth(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)

	claims, okClaim := c.Locals(util.KeyUser).(jwt.MapClaims)
	if !okClaim {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", "Unauthorize"))
	}

	var id int
	idClaim, ok := claims["ID"]

	if ok {
		id, ok = util.ToInt(idClaim)
	}

	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(util.ToMapKey("message", "invalid user id"))
	}

	if queryParam.Filter == nil {
		queryParam.Filter = make(map[string]interface{})
		queryParam.Filter["user_id"] = uint(id)
	}

	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *UserBookHandler) BorrowBooks(c *fiber.Ctx) error {
	id, ok := util.ToInt(c.Params("id"))
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(util.ToMapKey("message", "invalid path params"))
	}

	response, err := handler.service.BorrowBook(uint(id), c.Body())
	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *UserBookHandler) ReturnBooks(c *fiber.Ctx) error {
	id, ok := util.ToInt(c.Params("id"))
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(util.ToMapKey("message", "invalid path params"))
	}

	response, err := handler.service.ReturnBook(uint(id), c.Body())
	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

var userBookHandler *UserBookHandler
var onceUserBookHandler sync.Once

func GetUserBookHandler(db *gorm.DB) *UserBookHandler {
	onceUserBookHandler.Do(func() {
		userBookHandler = &UserBookHandler{service.GetUserBookService(db)}
	})

	return userBookHandler
}
