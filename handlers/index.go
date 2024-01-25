package handlers

import (
	"luuhai48/short/views"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func render(c *fiber.Ctx, comp templ.Component) error {
	handler := fasthttpadaptor.NewFastHTTPHandler(templ.Handler(comp))
	handler(c.Context())
	return nil
}

func redirect(c *fiber.Ctx, path string, status ...int) error {
	if len(status) > 0 {
		c.Status(status[0])
	} else {
		c.Status(200)
	}

	if c.Get("Hx-Request") != "" {
		c.Set("HX-Redirect", path)
		return nil
	}

	return c.Redirect(path)
}

func HandleGetHomeIndex(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return render(c, views.Index())
	}
	return redirect(c, "/short")
}
