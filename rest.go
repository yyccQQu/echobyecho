package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Id    string `json:"id" form:"id"`
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}

func saveUser(c echo.Context) error {
	u := new(User)
	u.Name = c.FormValue("name")
	u.Email = c.FormValue("email")
	u.Id = "1"

	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	u := new(User)
	id := c.Param("id")
	u.Name = "will"
	u.Email = "will@will.com"
	u.Id = id
	return c.JSON(http.StatusOK, u)
}

func updateUser(c echo.Context) error {
	u := new(User)
	id := c.Param("id")
	u.Name = "willf"
	u.Email = "will@willf.com"
	u.Id = id

	return c.JSON(http.StatusOK, u)
}

func deleteUser(c echo.Context) error {
	u := new(User)
	id := c.Param("id")
	println(u)
	println(id)
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()
	e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)
	e.Logger.Fatal(e.Start(":1323"))
}