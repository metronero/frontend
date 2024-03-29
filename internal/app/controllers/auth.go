package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/metronero/metronero-go/models"

	"gitlab.com/metronero/frontend/internal/utils/config"
)

func GetLogin(c *fiber.Ctx) error {
	// Display a success message if the user ended up here
	// after successfully logging in.
	reg := c.Query("success", "false")

	// Display a message telling the user to login to proceed
	// in case of invalid or expired token.
	exp := c.Query("expired", "false")
	return c.Render("login", fiber.Map{
		"Registered": reg,
		"Expired":    exp,
	})
}

func PostLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "" || password == "" {
		return serveErrorPage(c, http.StatusBadRequest,
			"Required form fields must not be empty")
	}
	token, err := config.Api.PostLogin(username, password)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token.Token
	cookie.Expires = token.ValidUntil
	cookie.Secure = true
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	if username == "admin" {
		return c.Redirect("/admin/dashboard")
	}
	return c.Redirect("/merchant/dashboard")
}

func GetLogout(c *fiber.Ctx) error {
	c.ClearCookie("token")
	return c.Redirect("/login")
}

func AdminCreateMerchant(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	token := c.Cookies("token")
	if err := config.Api.PostAdminRegister(token, &models.NewAccount{
			Username: username, Password: password, Role: "shop",
		}); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/admin/merchants?created=true")
}
