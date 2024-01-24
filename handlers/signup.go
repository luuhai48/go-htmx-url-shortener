package handlers

import (
	"luuhai48/short/utils"
	"luuhai48/short/views/signup"

	"github.com/gofiber/fiber/v2"
)

func HandleGetSignupIndex(c *fiber.Ctx) error {
	return render(c, signup.Index())
}

func ValidatePostSignup(username string, password string, passwordConfirm string) (hasError bool, errors signup.SignupErrors) {
	if err := utils.ValidateUsernameFormat(username); err != nil {
		hasError = true
		errors.Username = err.Error()
	}
	if err := utils.ValidatePasswordFormat(password); err != nil {
		hasError = true
		errors.Password = err.Error()
	}
	if passwordConfirm != password {
		hasError = true
		errors.PasswordConfirm = "Confirm password doesn't match"
	}
	return hasError, errors
}

func HandlePostSignup(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	passwordConfirm := c.FormValue("passwordConfirm")

	hasError, errors := ValidatePostSignup(username, password, passwordConfirm)
	if !hasError {
		return render(c, signup.Success())
	}

	return render(c, signup.Form(signup.SignupParams{
		Username:        username,
		Password:        password,
		PasswordConfirm: passwordConfirm,
		HasError:        hasError,
		Errors:          errors,
	}))
}
