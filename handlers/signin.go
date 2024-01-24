package handlers

import (
	"errors"
	"net/url"
	"strings"
	"time"

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
	remember := c.FormValue("remember", "off")

	params := signin.SigninParams{
		Username: username,
		Password: password,
		Remember: remember,
		HasError: false,
		Errors:   signin.SigninErrors{},
	}

	user := ValidatePostSignin(&params)

	if !params.HasError {
		session := &models.Session{
			UserID:     user.ID,
			Username:   user.Username,
			ValidUntil: time.Now().Add(time.Hour * 24 * 30),
		}
		if err := models.CreateSession(session); err != nil {
			params.Errors.General = err.Error()
			return render(c, signin.Form(params))
		}

		cookie := &fiber.Cookie{
			Name:        "session",
			Value:       session.ID,
			Path:        "/",
			Secure:      true,
			HTTPOnly:    true,
			SessionOnly: true,
		}

		if remember == "on" {
			cookie.SessionOnly = false
			cookie.Expires = session.ValidUntil
		}

		c.Cookie(cookie)

		next := c.Query("next", "/")
		if c.Get("Hx-Current-Url") != "" {
			u, err := url.Parse(c.Get("Hx-Current-Url"))
			if err == nil {
				n := u.Query().Get("next")
				if n != "" {
					next = n
				}
			}
		}

		parsed, err := url.QueryUnescape(next)
		if err == nil {
			next = parsed
		}

		return redirect(c, next)
	}

	return render(c, signin.Form(params))
}
