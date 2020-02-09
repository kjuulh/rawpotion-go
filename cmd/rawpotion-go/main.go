package main

import (
	"fmt"
	"net/http"

	"github.com/kjuulh/rawpotion-go/pkg/database"
	"github.com/kjuulh/rawpotion-go/pkg/tables"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Load config
	d := database.NewDatabase()
	d.LoadConfigFromFile("../../configs/config.yml")
	d.Open()

	defer d.Close()

	t, err := tables.NewUsersTable(tables.UsersTableConfig{Db: &d})
	if err != nil {
		fmt.Println("Failed at create table")
	}
	u, err := t.InsertUser(tables.UsersRow{
		Username: "Kasper J. Hermansen",
		Password: "Blizzar1",
	})
	if err != nil {
		fmt.Println("failed at insert")
	}
	fmt.Printf("Id: %s, Username: %s, Password: %s", u.Id, u.Username, u.Password)

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
