package admin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jinzhu/copier"
	"github.com/valyala/bytebufferpool"
	"gowek/repo"
	"html/template"
	"log"
	"time"
)

type User struct {
	Login           string
	Email           string
	UserID          uint
	IsAuthenticated bool
	IsSuperuser     bool
}

var userSession *session.Store

func NewUserFromMap(m map[string]interface{}) (u User) {
	return User{
		Login:           m["login"].(string),
		UserID:          uint(m["id"].(float64)),
		IsAuthenticated: true,
		IsSuperuser:     m["admin"].(bool),
	}
}

func Init(app *fiber.App) {

	userSession = session.New()

	//перехватываем каждый запрос и вставляем login для всех шаблонов
	authMiddleware(app)

	//форма логина
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("admin/login", fiber.Map{})
	})

	//процедура логина из формы
	app.Post("/login", loginPost)
	app.Get("/logout", logout)

	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.Render("admin/admin", fiber.Map{})
	})
}

func authMiddleware(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		s, _ := userSession.Get(c)
		user := User{IsAuthenticated: false, IsSuperuser: false}

		if v := s.Get("user"); v != nil {
			user = v.(User)
		}

		// добавляем user в список переменных контекста и список переменных шалобнов
		c.Locals("User", user)
		c.Bind(fiber.Map{
			"User": user,
		})

		// переходим к следующему обрабочику в стеке (списке) обработчиков запроса
		return c.Next()
	})

}

func loginPost(c *fiber.Ctx) error {
	login := c.FormValue("login")
	password := c.FormValue("password")

	dbuser, ok := repo.GetUser(login, password)
	if !ok {
		return c.SendString("Login or password not valid") //.Render("error", fiber.Map{"Error": "Login not found"})
	}

	var user User
	//копируем структуру из данных БД в локальную структуру которую будем хранить в сессии
	err := copier.Copy(&user, &dbuser)
	if err != nil {
		log.Println(err)
		return c.SendString(err.Error())
	}
	user.IsAuthenticated = true

	// Get or create session
	s, _ := userSession.Get(c)
	// Get session ID
	sid := s.ID()
	s.Set("user", user)
	s.Set("sid", sid)
	s.Set("ip", c.Context().RemoteIP().String())
	s.Set("logintime", time.Unix(time.Now().Unix(), 0).UTC().String())
	s.Set("ua", string(c.Request().Header.UserAgent()))

	err = s.Save()
	if err != nil {
		log.Println(err)
		return c.SendString(err.Error())
	}

	//редирект на главную старинцу через заголовок htmx
	//т.к. если делать c.Redirect("/"), то htmx вставит всю главную страницу в элемент error вместо редиректа
	c.Set("HX-Redirect", "/")
	//return c.Redirect("/")
	return c.SendStatus(fiber.StatusOK)
}
func logout(c *fiber.Ctx) error {
	s, _ := userSession.Get(c)
	if !s.Fresh() {
		s.Destroy()
	}
	return c.Redirect("/")
}

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
