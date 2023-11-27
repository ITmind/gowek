package restapi

import (
	"gowek/auth"
	"gowek/repo"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {

	api := app.Group("/api")

	api.Get("/users", getAllUsers)
	api.Get("/users/:login", getUser)
	api.Post("/users", addUser)

	api.Get("/notes", getNotesByUser)
	api.Get("/notes/:id", getNote)
	api.Post("/notes", addNote)
}

func getAllUsers(c *fiber.Ctx) error {
	res, err := repo.GetAllUsers()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal DB error")
	}
	return c.JSON(res)
}

func getUser(c *fiber.Ctx) error {
	login := c.Params("id")
	res, err := repo.GetUser(login)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal DB error")
	}
	return c.JSON(res)
}

func getNote(c *fiber.Ctx) error {
	userData := c.Locals("user")
	if userData == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := userData.(auth.User)

	i, _ := strconv.Atoi(c.Params("id"))
	noteID := uint(i)
	res, err := repo.GetNote(user.UserID, noteID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal DB error")
	}
	return c.JSON(res)
}

func addNote(c *fiber.Ctx) error {
	userData := c.Locals("user")
	if userData == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := userData.(auth.User)

	obj := new(repo.Note)
	if err := c.BodyParser(obj); err != nil {
		return fiber.NewError(http.StatusBadRequest, "bad request")
	}
	obj.ID = user.UserID

	err := repo.AddNote(obj)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal DB error")
	}

	return c.JSON(obj)
}

func getNotesByUser(c *fiber.Ctx) error {
	userData := c.Locals("user")
	if userData == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user := userData.(auth.User)

	notes, err := repo.GetAllNotesByUser(user.UserID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal DB error")
	}

	return c.JSON(notes)
}

func addUser(c *fiber.Ctx) error {

	userData := c.Locals("user")
	if userData == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	obj := new(repo.User)
	if err := c.BodyParser(obj); err != nil {
		return fiber.NewError(http.StatusBadRequest, "bad request")
	}

	err := repo.AddUser(obj)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal DB error")
	}

	return c.SendStatus(fiber.StatusCreated)
}
