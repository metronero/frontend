package controllers

import (
	"net/http"
	"strconv"
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/moneropay/go-monero/walletrpc"
	"gitlab.com/moneropay/metronero/metronero-backend/app/models"

	"gitlab.com/moneropay/metronero/metronero-frontend/app/api"
)

func AdminDashboard(c *fiber.Ctx) error {
	token := c.Cookies("token")
	resp, err := api.GetAdminDashboard(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("admin-dashboard", fiber.Map{
		"PageTitle":   "Dashboard",
		"Balance":     walletrpc.XMRToDecimal(resp.Stats.WalletBalance),
		"Profits":     walletrpc.XMRToDecimal(resp.Stats.TotalProfits),
		"Withdrawals": resp.RecentWithdrawals,
	}, "layouts/admin-panel")
}

func AdminInstance(c *fiber.Ctx) error {
	success := c.Query("success", "false")
	failed := c.Query("failed", "false")
	missing := c.Query("required", "false")

	token := c.Cookies("token")
	resp, err := api.GetInstanceInfo(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("admin-instance", fiber.Map{
		"PageTitle": "Instance",
		"InstanceInfo": resp,
		"DefaultCommission": walletrpc.XMRToDecimal(resp.Details.DefaultCommission),
		"Success": success,
		"Failed": failed,
		"Missing": missing,
	}, "layouts/admin-panel")
}

func AdminEditInstance(c *fiber.Ctx) error {
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
		err error
	)


	if conf.CustodialMode, err = strconv.ParseBool(custodialMode); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	floatCommission, err := strconv.ParseFloat(defaultCommission, 64)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	conf.DefaultCommission = uint64(floatCommission * 1000000000000)

	if conf.RegistrationsAllowed, err = strconv.ParseBool(registrationsAllowed); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	conf.WithdrawalTimes = withdrawalTimes

	log.Printf("%v", conf)
	token := c.Cookies("token")
	if err := api.UpdateInstance(token, &conf); err != nil {
		return c.Redirect("/admin/instance?failed=true")
	}

	return c.Redirect("/admin/instance?success=true")
}

func AdminWithdrawals(c *fiber.Ctx) error {
	return c.Render("admin-withdrawals", fiber.Map{
		"PageTitle": "Withdrawals",
	}, "layouts/admin-panel")
}

func AdminPayments(c *fiber.Ctx) error {
	return c.Render("admin-payments", fiber.Map{
		"PageTitle": "Payments",
	}, "layouts/admin-panel")
}

func AdminMerchants(c *fiber.Ctx) error {
	return c.Render("admin-merchants", fiber.Map{
		"PageTitle": "Merchants",
	}, "layouts/admin-panel")
}

func AdminMerchantEdit(c *fiber.Ctx) error {
	params := c.AllParams()
	return c.Render("admin-merchant-edit", fiber.Map{
		"Username":  params["uname"],
		"PageTitle": "Merchant Edit",
	}, "layouts/admin-panel")
}
