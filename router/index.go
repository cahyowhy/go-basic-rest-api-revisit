package router

import (
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
	Handler     handler.HandlerRoute
	Middlewares []handler.Adapter
}

func GetDefinedRoutes(db *gorm.DB) []Router {
	userHandler := handler.GetUserHandler(db)

	return []Router{
		{Path: "/users", Method: GET, Middlewares: []handler.Adapter{middleware.ParseQueryFilter}, Handler: userHandler.GetAllUsers},
	}
}
