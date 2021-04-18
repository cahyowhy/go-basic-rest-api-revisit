package main

import (
	"fmt"
	"net/http"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"github.com/cahyowhy/go-basit-restapi-revisit/handler"
	"github.com/cahyowhy/go-basit-restapi-revisit/router"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type App struct {
	db     *gorm.DB
	Router *mux.Router
}

func (app *App) Initialize(paramConfig *config.Config) {
	util.InitLogger()
	app.db = database.GetDatabase(paramConfig)
	app.setRouter()
}

func (app *App) Run(host string) {
	fmt.Printf("Running on host %s", host)
	util.ErrorLogger.Fatal(http.ListenAndServe(host, app.Router))
}

func (app *App) setRouter() {
	app.Router = mux.NewRouter()
	subrouter := app.Router.PathPrefix("/api/").Subrouter()

	for _, route := range router.GetDefinedRoutes(app.db) {
		newRoute := route

		var addapt = handler.Adapt(newRoute.Handler, newRoute.Middlewares...).ServeHTTP

		subrouter.HandleFunc(newRoute.Path, addapt).Methods(string(newRoute.Method))
	}
}
