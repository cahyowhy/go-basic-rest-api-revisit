package app

import (
	"fmt"
	"sync"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/database"
	"github.com/cahyowhy/go-basit-restapi-revisit/handler"
	"github.com/cahyowhy/go-basit-restapi-revisit/middleware"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gorm.io/gorm"
)

const (
	POST   string = "POST"
	DELETE string = "DELETE"
	GET    string = "GET"
	PUT    string = "PUT"
)

type App struct {
	DB       *gorm.DB
	FiberApp *fiber.App
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
	app.DB = database.GetDatabase(paramConfig)
	app.FiberApp = fiber.New()

	util.InitLogger()
	app.applyRoutes()
}

func (app *App) Run(host string) {
	fmt.Printf("Running on host %s", host)
	util.ErrorLogger.Fatal(app.FiberApp.Listen(host))
}

func (app *App) applyRoutes() {
	userHandler := handler.GetUserHandler(app.DB)
	bookHandler := handler.GetBookHandler(app.DB)
	userBookHandler := handler.GetUserBookHandler(app.DB)
	userFineHistoryHandler := handler.GetUserFineHistoryHandler(app.DB)

	app.FiberApp.Get("/monitor", monitor.New())

	// users
	app.FiberApp.Get("/api/users", middleware.AuthenticateJWT, middleware.ParseQueryFilter, userHandler.GetAll)
	app.FiberApp.Get("/api/users/paging/count", middleware.AuthenticateJWT, middleware.ParseQueryFilter, userHandler.Count)
	app.FiberApp.Post("/api/users", userHandler.Create)
	app.FiberApp.Get("/api/users/:id", middleware.AuthenticateJWT, userHandler.Get)
	app.FiberApp.Put("/api/users/:id", middleware.AuthenticateJWT, userHandler.Update)
	app.FiberApp.Post("/api/users/auth/login", userHandler.Login)
	app.FiberApp.Get("/api/users/auth/logout", middleware.AuthenticateJWT, userHandler.Logout)
	app.FiberApp.Get("/api/users/auth/session", userHandler.Session)

	// books
	app.FiberApp.Get("/api/books", middleware.AuthenticateJWT, middleware.ParseQueryFilter, bookHandler.GetAll)
	app.FiberApp.Get("/api/books/paging/count", middleware.AuthenticateJWT, middleware.ParseQueryFilter, bookHandler.Count)
	app.FiberApp.Post("/api/books", middleware.AuthenticateJWT, bookHandler.Create)
	app.FiberApp.Get("/api/books/:id", middleware.AuthenticateJWT, bookHandler.Get)
	app.FiberApp.Put("/api/books/:id", middleware.AuthenticateJWT, bookHandler.Update)

	// user-books
	app.FiberApp.Get("/api/user-books", middleware.AuthenticateJWT, middleware.ParseQueryFilter, userBookHandler.GetAll)
	app.FiberApp.Get("/api/user-books/paging/count", middleware.AuthenticateJWT, middleware.ParseQueryFilter, userBookHandler.Count)
	app.FiberApp.Get("/api/user-books/me", middleware.AuthenticateJWT, middleware.ParseQueryFilter, userBookHandler.GetAllFromAuth)
	app.FiberApp.Post("/api/user-books/borrows/:id", middleware.AuthenticateJWT, middleware.AuthenticateAdmin, userBookHandler.BorrowBooks)
	app.FiberApp.Put("/api/user-books/returns/:id", middleware.AuthenticateJWT, middleware.AuthenticateAdmin, userBookHandler.ReturnBooks)

	// user-fine-histories
	app.FiberApp.Get("/api/user-fine-histories", middleware.AuthenticateJWT, middleware.AuthenticateAdmin, middleware.ParseQueryFilter, userFineHistoryHandler.GetAll)
	app.FiberApp.Get("/api/user-fine-histories/paging/count", middleware.AuthenticateJWT, middleware.AuthenticateAdmin, middleware.ParseQueryFilter, userFineHistoryHandler.Count)
	app.FiberApp.Put("/api/user-fine-histories/pay/:id", middleware.AuthenticateJWT, middleware.AuthenticateAdmin, userFineHistoryHandler.PayBookFine)
}
