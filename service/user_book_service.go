package service

import (
	"database/sql"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type UserBookService struct {
	db   *gorm.DB
	base *baseService
}

func (service *UserBookService) FindAll(offset int, limit int, filter map[string]interface{}) (map[string]interface{}, error) {
	userBooks := []model.UserBook{}

	tx := service.db.Omit("users.password").Joins("User")
	tx = tx.Joins("Book").Offset(offset).Limit(limit)

	if filter != nil {
		tx = tx.Where(filter)
	}

	if err := tx.Find(&userBooks).Error; err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", userBooks), nil
}

func (service *UserBookService) Count(filter map[string]interface{}) (map[string]interface{}, error) {
	var total int64
	if err := service.base.count(&total, &model.UserBook{}, filter); err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", total), nil
}

func countAlreadyBorrowFrom(userId uint, service *UserBookService, bookIds ...uint) (<-chan int64, <-chan error) {
	r := make(chan int64)
	rErr := make(chan error)

	go func() {
		defer close(r)
		defer close(rErr)

		var total int64

		filter := make(map[string]interface{})
		filter["book_id"] = bookIds
		filter["user_id"] = userId
		filter["return_date"] = nil

		if err := service.db.Model(&model.UserBook{}).Where(filter).Count(&total).Error; err != nil {
			util.ErrorLogger.Fatal(err.Error())
			rErr <- err
		}

		r <- total
	}()

	return r, rErr
}

func (service *UserBookService) getParsedUserBooks(userId uint, body io.Reader) ([]model.UserBook, model.UserBorrowBook, []uint, map[string]interface{}) {
	userBorrowBook := model.UserBorrowBook{}
	if err := service.base.decodeJson(&userBorrowBook, body); err != nil {
		return nil, userBorrowBook, nil, util.ToMapKey("message", err.Error())
	}

	if err := config.GetValidator().Struct(userBorrowBook); err != nil {
		return nil, userBorrowBook, nil, util.ToMapKey("errors", util.ValidationErrToString(err.(validator.ValidationErrors)))
	}

	userBooks := []model.UserBook{}
	userBookIds := []uint{}

	for _, value := range userBorrowBook.Books {
		userBooks = append(userBooks, model.UserBook{BookID: value.BookId, UserID: userId, BorrowDate: time.Now()})
		userBookIds = append(userBookIds, value.BookId)
	}

	return userBooks, userBorrowBook, userBookIds, nil
}

func (service *UserBookService) ReturnBook(userId uint, body io.Reader) (map[string]interface{}, error) {
	userBooks, _, userBookBookIds, errMap := service.getParsedUserBooks(userId, body)
	if errMap != nil {
		return errMap, errors.New("error parsed from body")
	}

	var fine int
	err := service.db.Transaction(func(tx *gorm.DB) error {
		paramUserBook := make(map[string]interface{})
		paramUserBook["user_id"] = userId
		paramUserBook["book_id"] = userBookBookIds
		paramUserBook["return_date"] = nil

		if err := tx.Where(paramUserBook).Find(&userBooks).Error; err != nil {
			return err
		}

		userBookIds := []int32{}

		for _, value := range userBooks {
			fine += util.CountTotalFine(value.BorrowDate)
			if value.ID > 0 {
				userBookIds = append(userBookIds, int32(value.ID))
			}
		}

		if fine > 0 {
			if err := tx.Create(&model.UserFineHistory{UserID: userId, Fine: uint(fine), UserBookIds: userBookIds}).Error; err != nil {
				return err
			}
		}

		paramUpdate := make(map[string]interface{})
		paramUpdate["return_date"] = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		if err := tx.Model(&model.UserBook{}).Where(paramUserBook).Updates(paramUpdate).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	response := util.ToMapKey("message", "return book succeed")
	response["data"] = map[string]interface{}{"fine": fine}

	return response, nil
}

func (service *UserBookService) BorrowBook(userId uint, body io.Reader) (map[string]interface{}, error) {
	userBooks, _, userBookIds, errMap := service.getParsedUserBooks(userId, body)
	if errMap != nil {
		return errMap, errors.New("error parsed from body")
	}

	hasBorrowCh, errHasBorrowCh := countAlreadyBorrowFrom(userId, service, userBookIds...)
	totalBorrowCh, errTotalBorrowCh := countAlreadyBorrowFrom(userId, service)

	hasBorrow, errHasBorrow := <-hasBorrowCh, <-errHasBorrowCh
	totalBorrow, errTotalBorrow := <-totalBorrowCh, <-errTotalBorrowCh

	if errHasBorrow != nil {
		return util.ToMapKey("message", errHasBorrow.Error()), errHasBorrow
	} else if errTotalBorrow != nil {
		return util.ToMapKey("message", errTotalBorrow.Error()), errTotalBorrow
	} else if hasBorrow > 0 {
		err := errors.New("already borrow one of this book")

		return util.ToMapKey("message", err.Error()), err
	} else if totalBorrow >= 4 {
		err := errors.New("cannot borrow anymore, already borrow 4 book")

		return util.ToMapKey("message", err.Error()), err
	}

	if err := service.db.Create(&userBooks).Error; err != nil {
		return util.ToMapKey("message", err.Error()), err
	}

	return util.ToMapKey("data", userBooks), nil
}

var userBookService *UserBookService
var onceUserBookService sync.Once

func GetUserBookService(db *gorm.DB) *UserBookService {
	onceUserBookService.Do(func() {
		userBookService = &UserBookService{db, getBaseService(db)}
	})

	return userBookService
}
