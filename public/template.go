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
	return c.Render(http.StatusOK, "hello", "Will") //只能渲染数据，不能嵌套模板
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()

	// http://localhost:1323/static/index.html

	e.Static("/static","/Users/xieyadong/gopath1/src/echobyecho/public")

	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/hello", Hello)
	e.Logger.Fatal(e.Start(":1323"))
}
