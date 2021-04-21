package app

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"github.com/cahyowhy/go-basit-restapi-revisit/handler"
	"github.com/cahyowhy/go-basit-restapi-revisit/middleware"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Method string

const (
	POST   Method = "POST"
	DELETE Method = "DELETE"
	GET    Method = "GET"
	PUT    Method = "PUT"
)

type Router struct {
	Path        string
	Method      Method
	Handler     http.HandlerFunc
	Middlewares []handler.Adapter
}

type App struct {
	DB        *gorm.DB
	Router    *mux.Router
	Subrouter *mux.Router
}

var app *App
var onceApp sync.Once

func GetApp(config *config.Config) *App {
	onceApp.Do(func() {
		app = &App{}
		app.Initialize(config)
	})

	return app
}

func (app *App) Initialize(paramConfig *config.Config) {
	util.InitLogger()
	app.DB = database.GetDatabase(paramConfig)
	app.setRouter()
}

func (app *App) Run(host string) {
	fmt.Printf("Running on host %s", host)
	util.ErrorLogger.Fatal(http.ListenAndServe(host, app.Router))
}

func (app *App) setRouter() {
	app.Router = mux.NewRouter()
	app.Subrouter = app.Router.PathPrefix("/api/").Subrouter()
	ApplyRoutes(app)
}

func ApplyRoutes(app *App) {
	userHandler := handler.GetUserHandler(app.DB)
	bookHandler := handler.GetBookHandler(app.DB)
	userBookHandler := handler.GetUserBookHandler(app.DB)
	userFineHistoryHandler := handler.GetUserFineHistoryHandler(app.DB)

	routes := []Router{
		// users
		{Path: "/users", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.ParseQueryFilter}, Handler: userHandler.GetAll},
		{Path: "/users/paging/count", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.ParseQueryFilter}, Handler: userHandler.Count},
		{Path: "/users", Method: POST, Middlewares: []handler.Adapter{}, Handler: userHandler.Create},
		{Path: "/users/{id}", Method: PUT, Middlewares: []handler.Adapter{middleware.AuthenticateJWT}, Handler: userHandler.Update},
		{Path: "/users/{id}", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT}, Handler: userHandler.Get},
		{Path: "/users/auth/login", Method: POST, Middlewares: []handler.Adapter{}, Handler: userHandler.Login},
		{Path: "/users/auth/logout", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT}, Handler: userHandler.Logout},
		{Path: "/users/auth/session", Method: GET, Middlewares: []handler.Adapter{}, Handler: userHandler.Session},
		// books
		{Path: "/books", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.ParseQueryFilter}, Handler: bookHandler.GetAll},
		{Path: "/books/paging/count", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.ParseQueryFilter}, Handler: bookHandler.Count},
		{Path: "/books", Method: POST, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin}, Handler: bookHandler.Create},
		{Path: "/books/{id}", Method: PUT, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin}, Handler: bookHandler.Update},
		{Path: "/books/{id}", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT}, Handler: bookHandler.Get},
		// user-books
		{Path: "/user-books", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin, middleware.ParseQueryFilter}, Handler: userBookHandler.GetAll},
		{Path: "/user-books/me", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.ParseQueryFilter}, Handler: userBookHandler.GetAllFromAuth},
		{Path: "/user-books/paging/count", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin, middleware.ParseQueryFilter}, Handler: userBookHandler.Count},
		{Path: "/user-books/borrows/{id}", Method: POST, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin}, Handler: userBookHandler.BorrowBooks},
		{Path: "/user-books/returns/{id}", Method: PUT, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin}, Handler: userBookHandler.ReturnBooks},
		// user-fine-histories
		{Path: "/user-fine-histories", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin, middleware.ParseQueryFilter}, Handler: userFineHistoryHandler.GetAll},
		{Path: "/user-fine-histories/paging/count", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin, middleware.ParseQueryFilter}, Handler: userFineHistoryHandler.Count},
		{Path: "/user-fine-histories/pay/{id}", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT, middleware.AuthenticateAdmin, middleware.ParseQueryFilter}, Handler: userFineHistoryHandler.PayBookFine},
	}

	for _, route := range routes {
		newRoute := route
		var addapt = handler.Adapt(newRoute.Handler, newRoute.Middlewares...).ServeHTTP

		app.Subrouter.HandleFunc(newRoute.Path, addapt).Methods(string(newRoute.Method))
	}
}
