package main

import (
	"log"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"gorm.io/gorm"
)

func main() {
	db := database.GetDatabase(config.GetConfig())

	dbMigrate(dbDrop(db))
}

func dbDrop(db *gorm.DB) *gorm.DB {
	error := db.Migrator().DropTable(&model.UserFineHistory{}, &model.UserBook{}, &model.User{}, &model.Book{}, &model.UserSession{})

	if error != nil {
		log.Fatal("Failed drop")
	}

	return db
}

func dbMigrate(db *gorm.DB) *gorm.DB {
	error := db.AutoMigrate(&model.UserFineHistory{}, &model.UserBook{}, &model.User{}, &model.Book{}, &model.UserSession{})

	if error != nil {
		log.Fatal("Failed Migration")
	}

	return db
}
