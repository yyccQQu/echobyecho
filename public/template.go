package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

// 实现 echo.Renderer 接口
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "Will")
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	// 访问URI：/js/main.js 会寻找文件 assets/js/main.js
	e.Static("/", "assets")

	e.File("/", "public/index.html")

	e.Renderer = t
	e.GET("/hello", Hello)
	e.Logger.Fatal(e.Start(":1323"))
}
