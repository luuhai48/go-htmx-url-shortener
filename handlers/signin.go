package handlers

import (
	"errors"
	"log/slog"
	"strings"

	"luuhai48/short/models"
	"luuhai48/short/utils"
	"luuhai48/short/views/signin"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HandleGetSigninIndex(c *fiber.Ctx) error {
	return render(c, signin.Index())
}

func ValidatePostSignin(params *signin.SigninParams) *models.User {
	user, err := models.FindUserByUsername(params.Username)
	if err != nil {
		params.HasError = true
		if errors.Is(err, gorm.ErrRecordNotFound) {
			params.Errors.Username = "Account not registered!"
		} else {
			params.Errors.Username = err.Error()
		}
		return nil
	}

	if user.AccountStatus == models.AccountStateBlocked {
		params.HasError = true
		params.Errors.Username = "Account has been blocked!"
		return nil
	}

	if _, err := utils.VerifyPasswordHash(user.Password, params.Password); err != nil {
		params.HasError = true
		params.Errors.Password = "Invalid credentials!"
		return nil
	}

	return user
}

func HandlePostSignin(c *fiber.Ctx) error {
	username := strings.Trim(strings.ToLower(c.FormValue("username")), " ")
	password := c.FormValue("password")

	params := signin.SigninParams{
		Username: username,
		Password: password,
		HasError: false,
		Errors:   signin.SigninErrors{},
	}

	user := ValidatePostSignin(&params)

	if !params.HasError {
		slog.Info(user.ID)
		return redirect(c, "/")
	}

	return render(c, signin.Form(params))
}
