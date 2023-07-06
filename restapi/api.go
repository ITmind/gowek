package restapi

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gowek/admin"
	"gowek/repo"
	"net/http"
	"strconv"
)

func Init(app *fiber.App) {

	setJWTAuth(app)

	app.Get("/docs/*", swagger.HandlerDefault)

	api := app.Group("/api")

	api.Post("/token", login)
	api.Get("/users", getEntities[repo.User])
	api.Get("/users/:id", getEntityByID[repo.User])
	api.Post("/users", addUser)

	api.Get("/notes", getNotesByUser)
	api.Get("/notes/:id", getEntityByID[repo.Note])
	api.Post("/notes", addNote)
}

func getEntities[T repo.User | repo.Note](c *fiber.Ctx) error {
	res := repo.GetEntities[T]()
	return c.JSON(res)
}

// getEntityByID is a function to get entities by id
// @Summary Get entities by id
// @Description Get entities by id
// @Tags Entities
// @Accept json
// @Produce json
// @Success 200 {int} any
// @Router /api/notes [get]
func getEntityByID[T repo.User | repo.Note](c *fiber.Ctx) error {
	i, _ := strconv.Atoi(c.Params("id"))
	id := uint(i)
	res := repo.GetEntity[T](id)
	return c.JSON(res)
}

func addNote(c *fiber.Ctx) error {
	userData := c.Locals("user")
	if userData == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := admin.NewUserFromMap(userData.(map[string]interface{}))

	obj := new(repo.Note)
	if err := c.BodyParser(obj); err != nil {
		return fiber.NewError(http.StatusBadRequest, "bad request")
	}
	obj.ID = user.UserID

	repo.AddEntity(obj)
	return c.JSON(obj)
}

func getNotesByUser(c *fiber.Ctx) error {
	userData := c.Locals("user")
	if userData == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user := getUserFromToken(userData.(*jwt.Token))
	notes := repo.GetEntitiesByUser[repo.Note](user.UserID)
	return c.JSON(notes)
}

func addUser(c *fiber.Ctx) error {

	obj := new(repo.User)
	if err := c.BodyParser(obj); err != nil {
		return fiber.NewError(http.StatusBadRequest, "bad request")
	}

	repo.AddEntity(obj)

	return c.SendStatus(fiber.StatusCreated)
}
