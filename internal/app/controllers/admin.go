package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/frontend/internal/utils/config"
)

func GetAdminDashboard(c *fiber.Ctx) error {
	token := c.Cookies("token")
	resp, err := config.Api.GetAdmin(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	for i, w := range resp.RecentPayments {
		resp.RecentPayments[i].AmountFloat = walletrpc.XMRToDecimal(w.Amount)
	}
	return c.Render("admin-dashboard", fiber.Map{
		"PageTitle": "Dashboard",
		"Balance":   walletrpc.XMRToDecimal(resp.Stats.WalletBalance),
		"Profits":   walletrpc.XMRToDecimal(resp.Stats.TotalProfits),
		"Payments":  resp.RecentPayments,
	}, "layouts/admin-panel")
}
