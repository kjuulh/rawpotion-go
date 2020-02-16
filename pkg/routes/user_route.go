package routes

import (
	"github.com/kjuulh/rawpotion-go/pkg/rest"
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo) {
	e.POST("/users", rest.CreateUserHandler)
	e.GET("/users", rest.GetUsersHandler)
}
