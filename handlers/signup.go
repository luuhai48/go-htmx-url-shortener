package handlers

import (
	"strings"

	"luuhai48/short/models"
	"luuhai48/short/utils"
	"luuhai48/short/views/signup"

	"github.com/gofiber/fiber/v2"
)

func HandleGetSignupIndex(c *fiber.Ctx) error {
	return render(c, signup.Index())
}

func ValidatePostSignup(params *signup.SignupParams) {
	if err := utils.ValidateUsernameFormat(params.Username); err != nil {
		params.HasError = true
		params.Errors.Username = err.Error()
	}
	if err := utils.ValidatePasswordFormat(params.Password); err != nil {
		params.HasError = true
		params.Errors.Password = err.Error()
	}
	if params.PasswordConfirm != params.Password {
		params.HasError = true
		params.Errors.PasswordConfirm = "Confirm password doesn't match"
	}

	usernameExisted, err := models.CheckUsernameExists(params.Username)
	if err != nil {
		params.HasError = true
		params.Errors.Username = err.Error()
		return
	}
	if usernameExisted {
		params.HasError = true
		params.Errors.Username = "Username duplicated!"
		return
	}

	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		params.Errors.Password = err.Error()
		return
	}
	params.Password = hashedPassword
}

func HandlePostSignup(c *fiber.Ctx) error {
	username := strings.Trim(strings.ToLower(c.FormValue("username")), " ")
	password := c.FormValue("password")
	passwordConfirm := c.FormValue("passwordConfirm")

	params := signup.SignupParams{
		Username:        username,
		Password:        password,
		PasswordConfirm: passwordConfirm,
		HasError:        false,
		Errors:          signup.SignupErrors{},
	}
	ValidatePostSignup(&params)
	if !params.HasError {
		models.CreateUser(&models.User{
			Username: params.Username,
			Password: params.Password,
		})

		return render(c, signup.Success())
	}

	return render(c, signup.Form(params))
}
