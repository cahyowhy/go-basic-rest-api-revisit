package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var onceDb sync.Once
var db *gorm.DB

func GetDatabase(paramConfig *config.Config) *gorm.DB {
	onceDb.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", paramConfig.DbConfig.Host,
			paramConfig.DbConfig.Username, paramConfig.DbConfig.Password,
			paramConfig.DbConfig.Name, paramConfig.DbConfig.Port)

		dbRes, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Recorder.LogMode(logger.Silent),
		})

		if err != nil {
			log.Fatal("Could not connect database")
		}

		db = dbRes
	})

	return db
}
