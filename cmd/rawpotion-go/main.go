package main

import (
	"github.com/kjuulh/rawpotion-go/pkg/database"
	"github.com/kjuulh/rawpotion-go/pkg/routes"
	"github.com/kjuulh/rawpotion-go/pkg/tables"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Load config
	d := database.NewDatabase()
	d.LoadConfigFromFile("configs/config.yml")
	d.Open()

	defer d.Close()

	tables.InitUsersTable(&d)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Prometheus
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Recover
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	// Request ID
	e.Use(middleware.RequestID())

	// Routes
	routes.Routes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8082"))
}
