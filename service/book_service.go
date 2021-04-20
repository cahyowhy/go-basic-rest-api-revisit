package service

import (
	"io"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type BookService struct {
	db   *gorm.DB
	base *baseService
}

func (bookService *BookService) FindAll(offset int, limit int, filter map[string]interface{}) (map[string]interface{}, error) {
	books := []model.Book{}
	if err := bookService.base.findAll(&books, offset, limit, filter); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", books), nil
}

func (bookService *BookService) Count(filter map[string]interface{}) (map[string]interface{}, error) {
	var total int64
	if err := bookService.base.count(&total, &model.Book{}, filter); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", total), nil
}

func (bookService *BookService) Find(id int) (map[string]interface{}, error) {
	book := model.Book{}

	if err := bookService.base.findWhere(&book, id); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", book), nil
}

func (bookService *BookService) Update(id int, body io.Reader) (map[string]interface{}, error) {
	book := model.Book{}
	if err := bookService.base.decodeJson(&book, body); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	book.Model.ID = uint(id)
	if err := bookService.base.update(&book, id); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", book), nil
}

func (bookService *BookService) Create(body io.Reader) (map[string]interface{}, error) {
	book := model.Book{}
	if err := bookService.base.decodeJson(&book, body); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	if err := config.GetValidator().Struct(book); err != nil {
		return util.ToMapKey("errors", util.ValidationErrToString(err.(validator.ValidationErrors))), err
	}

	if err := bookService.base.create(&book); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", book), nil
}

var bookService *BookService
var onceBookService sync.Once

func GetBookService(db *gorm.DB) *BookService {
	onceBookService.Do(func() {
		bookService = &BookService{db, getBaseService(db)}
	})

	return bookService
}
