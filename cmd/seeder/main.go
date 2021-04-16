package main

import (
	"log"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"github.com/cahyowhy/go-basit-restapi-revisit/fake"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"gorm.io/gorm"
)

func main() { //go run cmd/seeder/main.go
	StartSeed(database.GetDatabase(config.GetConfig()), 10)
}

func StartSeed(db *gorm.DB, total int) {
	userCh, bookCh := seedUsers(db, total), seedBooks(db, total)
	users, books := <-userCh, <-bookCh

	seedUserBooks(db, books, users)
}

func seedUsers(db *gorm.DB, total int) <-chan []model.User {
	r := make(chan []model.User)

	go func() {
		defer close(r)

		var users = fake.GetUsers(total)
		err := db.Debug().Create(&users).Error

		if err != nil {
			log.Fatal(err.Error())
		}

		r <- users
	}()

	return r
}

func seedBooks(db *gorm.DB, total int) <-chan []model.Book {
	r := make(chan []model.Book)

	go func() {
		var books = fake.GetBooks(total)
		// fmt.Println(books[0])
		err := db.Debug().Create(&books).Error

		if err != nil {
			log.Fatal(err.Error())
		}

		r <- books
	}()

	return r
}

func seedUserBooks(db *gorm.DB, books []model.Book, users []model.User) []model.UserBook {
	var userBooks = fake.GetUserBooks(books, users)
	err := db.Debug().Create(&userBooks).Error

	if err != nil {
		log.Fatal(err.Error())
	}

	return userBooks
}
