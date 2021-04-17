package database

import (
	"fmt"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
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
			util.ErrorLogger.Fatal("Could not connect database")
		}

		db = dbRes
	})

	return db
}
