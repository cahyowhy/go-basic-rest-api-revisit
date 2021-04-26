package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"github.com/cahyowhy/go-basit-restapi-revisit/fake"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"gorm.io/gorm"
)

//go run cmd/seeder/main.go 10 .test.env
//go run cmd/seeder/main.go 10
//go run cmd/seeder/main.go
func main() {
	total := 10
	envFile := []string{}

	fmt.Println(os.Args)

	if len(os.Args) >= 2 {
		totalArgs, valid := util.ToInt(os.Args[1])

		if !valid {
			log.Fatal("specify valid total data")
		} else {
			total = totalArgs
		}

		envFile = os.Args[2:]
	}

	cf := config.GetConfig(envFile...)
	db := database.GetDatabase(cf)

	StartSeed(db, total)
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
			util.ErrorLogger.Fatal(err.Error())
		}

		r <- users
	}()

	return r
}

func seedBooks(db *gorm.DB, total int) <-chan []model.Book {
	r := make(chan []model.Book)

	go func() {
		var books = fake.GetBooks(total)
		err := db.Debug().Create(&books).Error

		if err != nil {
			util.ErrorLogger.Fatal(err.Error())
		}

		r <- books
	}()

	return r
}

func seedUserBooks(db *gorm.DB, books []model.Book, users []model.User) []model.UserBook {
	var userBooks = fake.GetUserBooks(books, users)
	err := db.Debug().Create(&userBooks).Error

	if err != nil {
		util.ErrorLogger.Fatal(err.Error())
	}

	return userBooks
}
