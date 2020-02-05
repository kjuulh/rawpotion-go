package main

import (
	"net/http"

	"github.com/kjuulh/rawpotion-go/pkg/config"
	"github.com/kjuulh/rawpotion-go/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Load config
	cfg := config.GetConfigFromFile("configs/config.yml")
	d := database.NewDatabase(database.Config{
		Database: cfg.Database.Database,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	})
	d.OpenConnection()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":8082"))
}

// Handler
func hello(c echo.Context) error {
	type User struct {
		Name  string `json:"name" xml:"name"`
		Email string `json:"email" xml:"email"`
	}

	u := &User{
		Name:  "Kasper",
		Email: "hermansendev@gmail.com",
	}
	return c.JSON(http.StatusOK, u)
}
