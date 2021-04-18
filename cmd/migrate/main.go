package main

import (
	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
)

func main() {
	model.DbMigrate(database.GetDatabase(config.GetConfig()))
}
