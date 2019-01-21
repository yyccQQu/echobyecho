package holder

import (
	"github.com/labstack/echo"
	"net/http"
)

type UserPost struct {
	username string
	pwd string
}

// URI参数
// e.GET("/users/:id", getUser)
func GetUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}


// e.POST("/save", save)
func SaveUser(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	//avatar, err := c.FormFile("avatar") 	获取头像
	return c.String(http.StatusOK, "name:" + name + ", email:" + email)
}
