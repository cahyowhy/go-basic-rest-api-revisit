package model

import (
	"log"

	"gorm.io/gorm"
)

func DbMigrate(db *gorm.DB) *gorm.DB {
	error := db.AutoMigrate(&UserFineHistory{}, &UserBook{}, &User{}, &Book{}, &UserSession{})

	if error != nil {
		log.Fatal("Failed Migration")
	}

	return db
}
