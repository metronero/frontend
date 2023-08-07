package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/metronero/metronero-go/models"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/frontend/internal/utils/config"
)

func GetAdminInstance(c *fiber.Ctx) error {
	success := c.Query("success", "false")
	failed := c.Query("failed", "false")
	missing := c.Query("required", "false")

	token := c.Cookies("token")
	resp, err := config.Api.GetAdminInstance(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	resp.Details.DefaultCommissionFloat = walletrpc.XMRToDecimal(resp.Details.DefaultCommission)
	resp.Stats.WalletBalanceFloat = walletrpc.XMRToDecimal(resp.Stats.WalletBalance)
	resp.Stats.TotalProfitsFloat = walletrpc.XMRToDecimal(resp.Stats.TotalProfits)

	return c.Render("admin-instance", fiber.Map{
		"PageTitle":    "Instance",
		"InstanceInfo": resp,
		"Success":      success,
		"Failed":       failed,
		"Missing":      missing,
	}, "layouts/admin-panel")
}

func PostAdminInstance(c *fiber.Ctx) error {
	custodialMode := c.FormValue("custodial_mode")
	defaultCommission := c.FormValue("default_commission")
	registrationsAllowed := c.FormValue("registrations_allowed")
	withdrawalTimes := c.FormValue("withdrawal_times")

	if custodialMode == "" || defaultCommission == "" ||
		registrationsAllowed == "" || withdrawalTimes == "" {
		return c.Redirect("/admin/instance?required=true")
	}

	var (
		conf models.Instance
		err  error
	)

	if conf.CustodialMode, err = strconv.ParseBool(custodialMode); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	if conf.DefaultCommission, err = toAtomic(defaultCommission); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	if conf.RegistrationsAllowed, err = strconv.ParseBool(registrationsAllowed); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	conf.WithdrawalTimes = withdrawalTimes

	token := c.Cookies("token")
	if err := config.Api.PostAdminInstance(token, &conf); err != nil {
		return c.Redirect("/admin/instance?failed=true")
	}

	return c.Redirect("/admin/instance?success=true")
}
