package controllers

import (
	"gitlab.com/metronero/frontend/utils/token"

	"github.com/gofiber/fiber/v2"
)

func serveErrorPage(c *fiber.Ctx, code int, message string) error {
	dUrl := "/admin/dashboard"
	if token.GetUsername(c) != "admin" {
		dUrl = "/merchant/dashboard"
	}
	return c.Render("error", fiber.Map{
		"Code": code,
		"Message": message,
		"DashboardUrl": dUrl,
	})
}
