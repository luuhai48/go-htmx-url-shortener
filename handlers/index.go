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

func HandleGetHomeIndex(c *fiber.Ctx) error {
	return render(c, views.Index())
}
