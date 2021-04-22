package handler

import (
	"net/http"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserFineHistoryHandler struct {
	service *service.UserFineHistoryService
}

func (handler *UserFineHistoryHandler) PayBookFine(c *fiber.Ctx) error {
	id, ok := util.ToInt(c.Params("id"))
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(util.ToMapKey("message", "invalid path params"))
	}

	if err := handler.service.PayBookFine(id, c.Body()); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.ToMapKey("message", err.Error()))
	}

	return c.JSON(util.ToMapKey("message", "Success paid fine"))
}

func (handler *UserFineHistoryHandler) GetAll(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)
	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *UserFineHistoryHandler) Count(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)
	response, err := handler.service.Count(queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

var userFineHistoryHandler *UserFineHistoryHandler
var onceUserFineHistoryHandler sync.Once

func GetUserFineHistoryHandler(db *gorm.DB) *UserFineHistoryHandler {
	onceUserFineHistoryHandler.Do(func() {
		userFineHistoryHandler = &UserFineHistoryHandler{service.GetUserFineHistoryService(db)}
	})

	return userFineHistoryHandler
}
