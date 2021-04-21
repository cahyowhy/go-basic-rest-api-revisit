package router

import (
	"net/http"

	"github.com/cahyowhy/go-basit-restapi-revisit/handler"
	"github.com/cahyowhy/go-basit-restapi-revisit/middleware"
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

func GetDefinedRoutes(db *gorm.DB) []Router {
	userHandler := handler.GetUserHandler(db)
	bookHandler := handler.GetBookHandler(db)
	userBookHandler := handler.GetUserBookHandler(db)
	userFineHistoryHandler := handler.GetUserFineHistoryHandler(db)

	return []Router{
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
}
