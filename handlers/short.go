package handlers

import (
	"log/slog"
	"net/url"

	"luuhai48/short/models"
	"luuhai48/short/utils"
	"luuhai48/short/views/short"

	"github.com/gofiber/fiber/v2"
)

func HandleShortIndex(c *fiber.Ctx) error {
	shorts, err := models.ListShortOfUser(c.Locals("user").(fiber.Map)["id"].(string))
	if err != nil {
		slog.Error(err.Error())
		return redirect(c, "/")
	}
	return render(c, short.Index(shorts))
}

func HandleNewShortIndex(c *fiber.Ctx) error {
	return render(c, short.NewShort())
}

func ValidatePostNewShort(params *short.ShortParams) {
	if _, err := url.ParseRequestURI(params.Url); err != nil {
		params.Errors.Url = err.Error()
		params.HasError = true
	}
}

func HandlePostNewShort(c *fiber.Ctx) error {
	params := short.ShortParams{
		Url:    c.FormValue("url"),
		Errors: short.ShortErrors{},
	}

	if !params.HasError {
		ID, _ := utils.GenShortID()

		new := &models.Short{
			BaseModel: models.BaseModel{
				ID: ID,
			},
			UserID: c.Locals("user").(fiber.Map)["id"].(string),
			Url:    params.Url,
		}

		if err := models.CreateShort(new); err != nil {
			params.HasError = true
			params.Url = err.Error()

			return render(c, short.NewForm(params))
		}

		return redirect(c, "/short")
	}

	return render(c, short.NewForm(params))
}

func HandleDeleteShort(c *fiber.Ctx) error {
	ID := c.FormValue("id")

	found, err := models.FindShortByID(ID)
	if err == nil && found.UserID == c.Locals("user").(fiber.Map)["id"].(string) {
		models.DeleteShortByID(ID)
	}

	return redirect(c, "/short")
}
