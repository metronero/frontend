package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/frontend/internal/utils/config"
	"gitlab.com/metronero/frontend/internal/utils/token"
)

func GetMerchantWithdrawals(c *fiber.Ctx) error {
	t := c.Cookies("token")
	resp, err := config.Api.GetMerchantWithdrawal(t)
	for i, w := range resp {
		resp[i].AmountFloat = walletrpc.XMRToDecimal(w.Amount)
	}
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("merchant-withdrawals", fiber.Map{
		"Username":    token.GetUsername(c),
		"PageTitle":   "Withdrawals",
		"Withdrawals": resp,
	}, "layouts/merchant-panel")
}

func GetAdminWithdrawals(c *fiber.Ctx) error {
	token := c.Cookies("token")
	resp, err := config.Api.GetAdminWithdrawal(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	for i, w := range resp {
		resp[i].AmountFloat = walletrpc.XMRToDecimal(w.Amount)
	}
	return c.Render("admin-withdrawals", fiber.Map{
		"PageTitle":   "Withdrawals",
		"Withdrawals": resp,
	}, "layouts/admin-panel")
}
