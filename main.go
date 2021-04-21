package main

import (
	"github.com/cahyowhy/go-basit-restapi-revisit/app"
	"github.com/cahyowhy/go-basit-restapi-revisit/config"
)

func main() {
	var configApp = config.GetConfig()
	app.GetApp(configApp).Run(":" + configApp.PORT)
}
