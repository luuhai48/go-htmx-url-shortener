package handlers

import (
	"luuhai48/short/models"

	"github.com/gofiber/fiber/v2"
)

func HandleGetSignOut(c *fiber.Ctx) error {
	cookie := c.Cookies("session")

	models.DeleteSessionById(cookie)
	c.ClearCookie("session")

	return redirect(c, "/")
}
