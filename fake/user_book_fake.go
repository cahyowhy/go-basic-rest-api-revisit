package fake

import (
	"database/sql"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/model"
)

func GetUserBooks(books []model.Book, users []model.User) []model.UserBook {
	var userBooks []model.UserBook

	for _, user := range users {
		for _, book := range books {
			var userBook = model.UserBook{
				UserID:     user.ID,
				BookID:     book.ID,
				BorrowDate: time.Date(2018, time.November, 12, 0, 0, 0, 0, time.UTC),
				ReturnDate: sql.NullTime{Time: time.Date(2020, time.November, 12, 0, 0, 0, 0, time.UTC), Valid: true},
			}

			userBooks = append(userBooks, userBook)
		}
	}

	return userBooks
}
