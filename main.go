package main

import (
	"github.com/cahyowhy/go-basit-restapi-revisit/config"
)

func main() {
	var configApp = config.GetConfig()

	app := App{}
	app.Initialize(configApp)
	app.Run(":" + configApp.PORT)
}
