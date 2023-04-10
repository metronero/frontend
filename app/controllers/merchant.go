package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/frontend/utils/config"
	"gitlab.com/metronero/frontend/utils/token"
)

func MerchantDashboard(c *fiber.Ctx) error {
	t := c.Cookies("token")
	resp, err := config.Api.GetMerchantInfo(t)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	resp.Stats.BalanceFloat = walletrpc.XMRToDecimal(resp.Stats.Balance)
	resp.Stats.TotalSalesFloat = walletrpc.XMRToDecimal(resp.Stats.TotalSales)
	for i, p := range resp.RecentPayments {
		resp.RecentPayments[i].AmountFloat = walletrpc.XMRToDecimal(p.Amount)
		resp.RecentPayments[i].FeeFloat = walletrpc.XMRToDecimal(p.Fee)
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
	for i, p := range resp {
		resp[i].AmountFloat = walletrpc.XMRToDecimal(p.Amount)
		resp[i].FeeFloat = walletrpc.XMRToDecimal(p.Fee)
	}
	return c.Render("merchant-payments", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Payments",
		"Payments": resp,
	}, "layouts/merchant-panel")
}

func MerchantWithdrawals(c *fiber.Ctx) error {
	t := c.Cookies("token")
	resp, err := config.Api.GetMerchantWithdrawals(t)
	for i, w := range resp {
		resp[i].AmountFloat = walletrpc.XMRToDecimal(w.Amount)
	}
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("merchant-withdrawals", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Withdrawals",
		"Withdrawals": resp,
	}, "layouts/merchant-panel")
}

func MerchantTemplate(c *fiber.Ctx) error {
	saved := c.Query("saved", "false")
	reset := c.Query("reset", "false")

	return c.Render("merchant-theme", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Theme",
		"Saved": saved,
		"Reset": reset,
	}, "layouts/merchant-panel")
}

func MerchantTemplatePreview(c *fiber.Ctx) error {
	t := c.Cookies("token")
	content, err := config.Api.MerchantTemplatePreview(t)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send(content)
}

func MerchantTemplateUpload(c *fiber.Ctx) error {
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
	if err := config.Api.MerchantTemplateUpload(t, f); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/merchant/template?saved=true")
}

func MerchantTemplateReset(c *fiber.Ctx) error {
	t := c.Cookies("token")
	if err := config.Api.MerchantTemplateReset(t); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/merchant/template?reset=true")
}

func MerchantAccount(c *fiber.Ctx) error {
	return c.Render("merchant-account", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Account",
	}, "layouts/merchant-panel")
}
