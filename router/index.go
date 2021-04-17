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

	return []Router{
		{Path: "/users", Method: GET, Middlewares: []handler.Adapter{middleware.ParseQueryFilter, middleware.AuthenticateJWT}, Handler: userHandler.GetAll},
		{Path: "/users", Method: POST, Middlewares: []handler.Adapter{}, Handler: userHandler.Create},
		{Path: "/users/{id}", Method: PUT, Middlewares: []handler.Adapter{middleware.AuthenticateJWT}, Handler: userHandler.Update},
		{Path: "/users/{id}", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT}, Handler: userHandler.Get},
		{Path: "/users/auth/login", Method: POST, Middlewares: []handler.Adapter{}, Handler: userHandler.Login},
		{Path: "/users/auth/logout", Method: GET, Middlewares: []handler.Adapter{middleware.AuthenticateJWT}, Handler: userHandler.Logout},
		{Path: "/users/auth/session", Method: GET, Middlewares: []handler.Adapter{}, Handler: userHandler.Session},
	}
}
