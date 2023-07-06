package main

import (
	"encoding/gob"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"gowek/admin"
	_ "gowek/docs"
	"gowek/repo"
	"gowek/restapi"
	"log"
	"os"
)

func init() {
	gob.Register(admin.User{})

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error load env file. Err: %s", err)
	}
}

func main() {

	if len(os.Args) > 1 {

		switch command := os.Args[1]; command {
		case "migrate":
			repo.Migrate()
		default:
			println("not valid arguments")
		}
	} else {
		app := NewApp()
		log.Fatal(app.Listen(":1323"))
	}

}

// NewApp Создает и возвращает fiber.App
// сделано для упрощения тестирования
func NewApp() *fiber.App {

	repo.Init()
	templateEngine := html.New("templates", ".go.html")

	app := fiber.New(fiber.Config{
		Views:       templateEngine,
		ViewsLayout: "base",
	})

	//систему аутентификации инициализируем первой, что бы ее middleware отрабатывало первой в стеке!
	admin.Init(app)
	restapi.Init(app)
	app.Use(recover.New())

	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	return app

}
