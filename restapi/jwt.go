package restapi

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gowek/admin"
	"gowek/repo"
	"os"
	"time"
)

var secret string

func setJWTAuth(app *fiber.App) {
	secret = os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "5HB0hBQ3YhbasmG1EL1yQxukG15AK1Pdgcez3ekfB"
	}

	app.Use("/api", jwtware.New(jwtware.Config{
		Filter: func(ctx *fiber.Ctx) bool {
			return ctx.OriginalURL() == "/api/token"
		},
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	}))

}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	dbuser, ok := repo.GetUser(user, pass)
	// Throws Unauthorized error
	if !ok {
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
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func getUserFromToken(token *jwt.Token) admin.User {
	m := token.Claims.(jwt.MapClaims)
	return admin.NewUserFromMap(m)
}
