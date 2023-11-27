package admin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"
	"html/template"
)

// RenderTemplate - рендер не через fiber
func RenderTemplate(c *fiber.Ctx) error {

	var data any = nil
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	//tmpl := make(map[string]*template.Template)

	fmap := template.FuncMap{
		"embed": func() string {
			return "test2"
		},
	}
	t, err := template.New("login.go.html").Funcs(fmap).ParseFiles("templates/login.go.html")

	if err != nil {
		println(err)
	}

	if err := t.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}

	// Set Content-Type to text/html
	c.Type("html", "UTF8")
	// Set rendered template to body
	c.Response().SetBody(buf.Bytes())

	return nil
}
