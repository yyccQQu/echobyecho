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
	return c.Render(http.StatusOK, "hello",  map[string]interface{}{
		"name": "Dolly!",
		"pwd": "123456",
	})
}
func Templatehtml(c echo.Context) error {
	//模板嵌套时相同key 会被覆盖 ---> 其他模板如果 没有该模板变量则不会被渲染
	return c.Render(http.StatusOK, "template",  map[string]interface{}{
		"names": "Dolly!",
		"pwd": "123456",
	})
}


// http://go-echo.org/guide/templates/
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
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
	e.GET("/template", Templatehtml)
	e.Logger.Fatal(e.Start(":1323"))
}
