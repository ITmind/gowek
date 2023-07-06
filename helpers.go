package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"
	"html/template"
	"strconv"
)

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}

func RenderTemplate(c *fiber.Ctx) error {

	var data any = nil
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	tmpl := make(map[string]*template.Template)

	fmap := template.FuncMap{
		"embed": func() string {
			return "test own render"
		},
	}
	t, err := template.New("base.go.html").Funcs(fmap).ParseFiles("templates/base.go.html")
	//t, err := template.New("test").Funcs(fmap).Parse("<h1>{{embed}}</h1>")
	if err != nil {
		println(err)
	}
	tmpl["index.html"] = template.Must(t, err)
	//tmpl["other.html"] = template.Must(template.ParseFiles("other.html", "base.html"))

	if err := t.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}

	// Set Content-Type to text/html
	c.Type("html", "UTF8")
	// Set rendered template to body
	c.Response().SetBody(buf.Bytes())

	return nil
}
