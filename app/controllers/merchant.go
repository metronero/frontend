package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"gitlab.com/metronero/frontend/utils/config"
	"gitlab.com/metronero/frontend/utils/token"
)

func MerchantDashboard(c *fiber.Ctx) error {
	t := c.Cookies("token")
	resp, err := config.Api.GetMerchantInfo(t)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("merchant-dashboard", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Dashboard",
		"Merchant":  resp,
	}, "layouts/merchant-panel")
}

func MerchantPayments(c *fiber.Ctx) error {
	t := c.Cookies("token")
	resp, err := config.Api.GetMerchantPayments(t)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("merchant-payments", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Payments",
		"Payments": resp,
	}, "layouts/merchant-panel")
}

func MerchantRequestPayment(c *fiber.Ctx) error {
	return c.Render("merchant-request-payment", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Request Payment",
	}, "layouts/merchant-panel")
}

func MerchantWithdrawals(c *fiber.Ctx) error {
	t := c.Cookies("token")
	resp, err := config.Api.GetMerchantWithdrawals(t)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("merchant-withdrawals", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Withdrawals",
		"Withdrawals": resp,
	}, "layouts/merchant-panel")
}

func MerchantTheme(c *fiber.Ctx) error {
	return c.Render("merchant-theme", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Theme",
	}, "layouts/merchant-panel")
}

func MerchantAccount(c *fiber.Ctx) error {
	return c.Render("merchant-account", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Account",
	}, "layouts/merchant-panel")
}
