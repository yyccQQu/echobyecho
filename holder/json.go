package holder

import (
	"net/http"

	"github.com/labstack/echo"
	"echobyecho/holder"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}



// /show?team=x-men&member=wolverine
func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:" + team + ", member:" + member)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		u := new(User)
		u.Name = "will"
		u.Email = "will@will.com"
		//返回 {"name": "will", "email": "will@will.com"}
		return c.JSON(http.StatusOK, u)
		//return c.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
	})

	e.GET("/show", show)

	// 路由
	e.POST("/save", holder.SaveUser)
	e.GET("/users/:id", holder.GetUser)
	//e.PUT("/users/:id", updateUser)
	//e.DELETE("/users/:id", deleteUser)


	e.Logger.Fatal(e.Start(":1323"))
}










