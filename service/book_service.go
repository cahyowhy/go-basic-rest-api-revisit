package service

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type BookService struct {
	db *gorm.DB
}

func (bookService *BookService) FindAll(offset int, limit int, filter map[string]interface{}) (map[string]interface{}, error) {
	books := []model.Book{}
	tx := bookService.db.Offset(offset).Limit(limit)

	if filter != nil {
		tx = tx.Where(filter)
	}

	if err := tx.Find(&books).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var body = make(map[string]interface{})
	body["data"] = books

	return body, nil
}

func (bookService *BookService) Find(id int) (map[string]interface{}, error) {
	book := model.Book{}

	if err := bookService.db.First(&book, id).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var body = make(map[string]interface{})
	body["data"] = book

	return body, nil
}

func (bookService *BookService) Update(id int, body io.Reader) (map[string]interface{}, error) {
	book := model.Book{}
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&book); err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	book.Model.ID = uint(id)

	if err := bookService.db.Updates(&book).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var response = make(map[string]interface{})
	response["data"] = book

	return response, nil
}

func (bookService *BookService) Create(body io.Reader) (map[string]interface{}, error) {
	book := model.Book{}
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&book); err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var err = config.GetValidator().Struct(book)
	if err != nil {
		response := make(map[string]interface{})
		response["errors"] = util.ValidationErrToString(err.(validator.ValidationErrors))

		return response, err
	}

	if err := bookService.db.Save(&book).Error; err != nil {
		return util.GetReponseMessage(err.Error()), err
	}

	var response = make(map[string]interface{})
	response["data"] = book

	return response, nil
}

var bookService *BookService
var onceBookService sync.Once

func GetBookService(db *gorm.DB) *BookService {
	onceBookService.Do(func() {
		bookService = &BookService{db}
	})

	return bookService
}
