package auth

import (
	"gowek/repo"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jinzhu/copier"

	"log/slog"
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

func userFromMap(m map[string]interface{}) (u User) {
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

	setJWTAuth(app, "/api")
	app.Post("/token", getToken)
}

func authMiddleware(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		s, _ := userSession.Get(c)
		user := User{IsAuthenticated: false, IsSuperuser: false}

		if v := s.Get("user"); v != nil {
			user = v.(User)
		}

		// добавляем user в список переменных контекста и список переменных шалобнов
		c.Locals("user", user)
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

	dbuser, err := repo.GetUser(login)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal DB error")
	}

	if dbuser.Hash != password {
		slog.Error("Login or password not valid")
		return c.SendString("Login or password not valid")
	}

	var user User
	//копируем структуру из данных БД в локальную структуру которую будем хранить в сессии
	err = copier.Copy(&user, &dbuser)
	if err != nil {
		slog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
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
		slog.Error("save user session", "Err", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	//редирект на главную старинцу через заголовок htmx
	//т.к. если делать c.Redirect("/"), то htmx вставит всю главную страницу в элемент вместо редиректа
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
