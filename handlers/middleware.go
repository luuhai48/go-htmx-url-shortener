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
	if err != nil || !session.Valid || session.ValidUntil.Before(time.Now()) {
		c.ClearCookie("session")
		return redirect(c, signinPath)
	}

	c.Locals("user", fiber.Map{
		"username": session.Username,
		"id":       session.UserID,
	})

	return c.Next()
}

func OptionalAuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("session")
	if cookie != "" {
		session, err := models.FindSessionByID(cookie)
		if err != nil || !session.Valid || session.ValidUntil.Before(time.Now()) {
			c.ClearCookie("session")
			return c.Next()
		}

		c.Locals("user", fiber.Map{
			"username": session.Username,
			"id":       session.UserID,
		})
	}

	return c.Next()
}

func NoAuthGuard(c *fiber.Ctx) error {
	cookie := c.Cookies("session")
	if cookie == "" {
		return c.Next()
	}

	return redirect(c, "/")
}
