package rest

import (
	"net/http"

	"github.com/kjuulh/rawpotion-go/pkg/tables"
	"github.com/labstack/echo/v4"
)

type User struct {
	Id       string `json:"id" form:"id" query:"id"`
	Name     string `json:"name" form:"name" query:"name"`
	Password string `json:"password" form:"password" query:"password"`
}

func CreateUserHandler(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}

	user, err := tables.User.Insert(tables.UsersRow{
		Username: u.Name,
		Password: u.Password,
	})
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, user)
}

func GetUsersHandler(c echo.Context) (err error) {
	users, err := tables.User.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}
