package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"github.com/cahyowhy/go-basit-restapi-revisit/handler"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/router"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type App struct {
	db     *gorm.DB
	Router *mux.Router
}

func (app *App) Initialize(paramConfig *config.Config) {
	app.db = model.DbMigrate(database.GetDatabase(paramConfig))
	app.setRouter()
}

func (app *App) Run(host string) {
	fmt.Printf("Running on host %s", host)
	log.Fatal(http.ListenAndServe(host, app.Router))
}

func (app *App) setRouter() {
	app.Router = mux.NewRouter()
	subrouter := app.Router.PathPrefix("/api/").Subrouter()

	for _, route := range router.GetDefinedRoutes(app.db) {
		newRoute := route
		var routeHandler http.HandlerFunc = func(writer http.ResponseWriter, req *http.Request) {
			newRoute.Handler(app.db, writer, req)
		}

		var addapt = handler.Adapt(routeHandler, newRoute.Middlewares...).ServeHTTP

		subrouter.HandleFunc(newRoute.Path, addapt).Methods(string(newRoute.Method))
	}
}
