package handlers

import (
	"luuhai48/short/views/url"

	"github.com/gofiber/fiber/v2"
)

func HandleUrlHomeIndex(c *fiber.Ctx) error {
	return render(c, url.Index())
}
