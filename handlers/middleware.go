package handlers

import (
	"fmt"
	"luuhai48/short/models"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("session")
	signinPath := fmt.Sprintf("/signin?next=%s", url.QueryEscape(c.Path()))
	if cookie == "" {
		return redirect(c, signinPath)
	}

	session, err := models.FindSessionByID(cookie)
	if err != nil || session.ValidUntil.Before(time.Now()) {
		return redirect(c, signinPath)
	}

	return c.Next()
}
