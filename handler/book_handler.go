package handler

import (
	"net/http"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/service"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BookHandler struct {
	service *service.BookService
}

func (handler *BookHandler) GetAll(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)
	response, err := handler.service.FindAll(queryParam.Offset, queryParam.Limit, queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *BookHandler) Count(c *fiber.Ctx) error {
	queryParam := GetQueryParam(c)
	response, err := handler.service.Count(queryParam.Filter)

	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *BookHandler) Get(c *fiber.Ctx) error {
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

func (handler *BookHandler) Create(c *fiber.Ctx) error {
	response, err := handler.service.Create(c.Body())
	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(response)
}

func (handler *BookHandler) Update(c *fiber.Ctx) error {
	id, ok := util.ToInt(c.Params("id"))
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(util.ToMapKey("message", "invalid path params"))
	}

	response, err := handler.service.Update(int(id), c.Body())
	if err == nil {
		return c.JSON(response)
	}

	return c.Status(http.StatusInternalServerError).JSON(util.ToMapKey("message", err.Error()))
}

var bookHandler *BookHandler
var onceBookHandler sync.Once

func GetBookHandler(db *gorm.DB) *BookHandler {
	onceBookHandler.Do(func() {
		bookHandler = &BookHandler{service.GetBookService(db)}
	})

	return bookHandler
}
