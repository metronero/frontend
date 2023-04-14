package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/metronero/backend/pkg/models"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/frontend/internal/utils/config"
	"gitlab.com/metronero/frontend/internal/utils/token"
)

func GetMerchantDashboard(c *fiber.Ctx) error {
	t := c.Cookies("token")
	resp, err := config.Api.GetMerchant(t)
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

func AdminGetMerchant(c *fiber.Ctx) error {
	success := c.Query("success", "false")
	token := c.Cookies("token")
	merchant, err := config.Api.GetAdminMerchantById(token, c.Params("id"))
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("admin-merchant-edit", fiber.Map{
		"PageTitle":      "Merchant Edit",
		"Merchant":       merchant,
		"CommissionRate": walletrpc.XMRToDecimal(merchant.CommissionRate),
		"Success":        success,
	}, "layouts/admin-panel")
}

func GetMerchantAccount(c *fiber.Ctx) error {
	return c.Render("merchant-account", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Account",
	}, "layouts/merchant-panel")
}

func AdminMerchants(c *fiber.Ctx) error {
	token := c.Cookies("token")
	deleted := c.Query("deleted", "false")
	resp, err := config.Api.GetAdminMerchant(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("admin-merchants", fiber.Map{
		"PageTitle": "Merchants",
		"Merchants": resp,
		"Deleted":   deleted,
	}, "layouts/admin-panel")
}

func AdminEditMerchant(c *fiber.Ctx) error {
	commission := c.FormValue("commission")
	disable := c.FormValue("disable")

	if commission == "" && disable == "" {
		c.Redirect("/admin/merchants/" + c.Params("id"))
	}

	var conf models.MerchantSettings

	if commission != "" {
		floatCommission, err := strconv.ParseFloat(commission, 64)
		if err != nil {
			return serveErrorPage(c, http.StatusInternalServerError, err.Error())
		}
		comRate := uint64(floatCommission * 1000000000000)
		conf.CommissionRate = &comRate
	}

	if disable != "" {
		disabledBool, err := strconv.ParseBool(disable)
		if err != nil {
			return serveErrorPage(c, http.StatusInternalServerError, err.Error())
		}
		conf.Disabled = &disabledBool
	}

	id := c.Params("id")
	conf.AccountId = id
	token := c.Cookies("token")
	if err := config.Api.PostAdminMerchantById(token, id, conf); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect(fmt.Sprintf("/admin/merchants/%s?success=true", id))
}

func AdminDeleteMerchant(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if err := config.Api.DeleteAdminMerchantById(token, c.Params("id")); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/admin/merchants?deleted=true")
}
