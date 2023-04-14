package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"gitlab.com/metronero/frontend/internal/utils/config"
	"gitlab.com/metronero/frontend/internal/utils/token"
)

func GetMerchantTemplate(c *fiber.Ctx) error {
	saved := c.Query("saved", "false")
	reset := c.Query("reset", "false")

	return c.Render("merchant-theme", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Theme",
		"Saved":     saved,
		"Reset":     reset,
	}, "layouts/merchant-panel")
}

func GetMerchantTemplatePreview(c *fiber.Ctx) error {
	t := c.Cookies("token")
	content, err := config.Api.GetMerchantTemplate(t)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send(content)
}

func PostMerchantTemplate(c *fiber.Ctx) error {
	t := c.Cookies("token")
	file, err := c.FormFile("file")
	if err != nil {
		return serveErrorPage(c, http.StatusBadRequest, err.Error())
	}
	f, err := file.Open()
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	defer f.Close()
	if err := config.Api.PostMerchantTemplate(t, f); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/merchant/template?saved=true")
}

func PostMerchantTemplateReset(c *fiber.Ctx) error {
	t := c.Cookies("token")
	if err := config.Api.DeleteMerchantTemplate(t); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/merchant/template?reset=true")
}
