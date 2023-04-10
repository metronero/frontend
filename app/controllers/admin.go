package controllers

import (
	"net/http"
	"log"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/moneropay/go-monero/walletrpc"
	"gitlab.com/metronero/backend/app/models"

	"gitlab.com/metronero/frontend/utils/config"
)

func AdminDashboard(c *fiber.Ctx) error {
	token := c.Cookies("token")
	resp, err := config.Api.GetAdminDashboard(token)
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
	resp, err := config.Api.GetInstanceInfo(token)
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

	if conf.DefaultCommission, err = toAtomic(defaultCommission); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	if conf.RegistrationsAllowed, err = strconv.ParseBool(registrationsAllowed); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}

	conf.WithdrawalTimes = withdrawalTimes

	log.Printf("%v", conf)
	token := c.Cookies("token")
	if err := config.Api.UpdateInstance(token, &conf); err != nil {
		return c.Redirect("/admin/instance?failed=true")
	}

	return c.Redirect("/admin/instance?success=true")
}

func AdminWithdrawals(c *fiber.Ctx) error {
	token := c.Cookies("token")
	resp, err := config.Api.AdminGetWithdrawals(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("admin-withdrawals", fiber.Map{
		"PageTitle": "Withdrawals",
		"Withdrawals": resp,
	}, "layouts/admin-panel")
}

func AdminPayments(c *fiber.Ctx) error {
	token := c.Cookies("token")
	resp, err := config.Api.AdminGetPayments(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("admin-payments", fiber.Map{
		"PageTitle": "Payments",
		"Payments": resp,
	}, "layouts/admin-panel")
}

func AdminMerchants(c *fiber.Ctx) error {
	token := c.Cookies("token")
	deleted := c.Query("deleted", "false")
	resp, err := config.Api.GetMerchantList(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("admin-merchants", fiber.Map{
		"PageTitle": "Merchants",
		"Merchants": resp,
		"Deleted": deleted,
	}, "layouts/admin-panel")
}


func AdminGetMerchant(c *fiber.Ctx) error {
	success := c.Query("success", "false")
	token := c.Cookies("token")
	merchant, err := config.Api.GetMerchantById(token, c.Params("id"))
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Render("admin-merchant-edit", fiber.Map{
		"PageTitle": "Merchant Edit",
		"Merchant": merchant,
		"CommissionRate": walletrpc.XMRToDecimal(merchant.CommissionRate),
		"Success": success,
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
	if err := config.Api.AdminEditMerchant(token, id, conf); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect(fmt.Sprintf("/admin/merchants/%s?success=true", id))
}

func AdminDeleteMerchant(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if err := config.Api.DeleteMerchantById(token, c.Params("id")); err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/admin/merchants?deleted=true")
}
