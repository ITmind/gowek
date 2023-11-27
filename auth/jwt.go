package auth

import (
	"gowek/repo"
	"log/slog"
	"net/http"
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var secret string

func setJWTAuth(app *fiber.App, path string) {
	secret = os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "5HB0hBQ3YhbasmG1EL1yQxukG15AK1Pdgcez3ekfB"
	}

	app.Use(path, jwtware.New(jwtware.Config{
		ContextKey: "jwtuser",
		Filter: func(ctx *fiber.Ctx) bool {
			return ctx.OriginalURL() == "/api/token"
		},
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		SuccessHandler: func(c *fiber.Ctx) error {
			c.Locals("user", getUserFromToken(c.Locals("jwtuser").(*jwt.Token)))
			return c.Next()
		},
	}))

}

func getToken(c *fiber.Ctx) error {
	user := c.FormValue("username")
	pass := c.FormValue("password")

	dbuser, err := repo.GetUser(user)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal DB error")
	}
	// Throws Unauthorized error
	if dbuser.Hash != pass {
		slog.ErrorContext(c.Context(), "User not found", "username", user)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"login": dbuser.Login,
		"id":    dbuser.ID,
		"admin": false,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.ErrorContext(c.Context(), "Generate token error", "error", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func getUserFromToken(token *jwt.Token) User {
	m := token.Claims.(jwt.MapClaims)
	return userFromMap(m)
}
