package controllers

import (
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/metronero/backend/pkg/models"
	"gitlab.com/moneropay/go-monero/walletrpc"

	"gitlab.com/metronero/frontend/internal/utils/config"
	"gitlab.com/metronero/frontend/internal/utils/token"
)

func GetMerchantPayments(c *fiber.Ctx) error {
	t := c.Cookies("token")
	resp, err := config.Api.GetMerchantPayment(t)
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
		"Payments":  resp,
	}, "layouts/merchant-panel")
}

func GetPaymentPage(c *fiber.Ctx) error {
	resp, err := config.Api.GetPaymentPage(c.Params("id"))
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	c.Set("Refresh", "10")
	return c.Send(resp)
}

func GetMerchantRequest(c *fiber.Ctx) error {
	id := c.Query("id", "")
	created := c.Query("created", "false")

	var (
		reqUrl string
		err    error
	)
	if id != "" {
		reqUrl, err = url.JoinPath(config.Hostname, "p", id)
		if err != nil {
			return serveErrorPage(c, http.StatusInternalServerError, err.Error())
		}
	}
	return c.Render("merchant-request-payment", fiber.Map{
		"Username":      token.GetUsername(c),
		"PageTitle":     "Request Payment",
		"PaymentReqUrl": reqUrl,
		"Created":       created,
	}, "layouts/merchant-panel")
}

func PostMerchantRequest(c *fiber.Ctx) error {
	amount := c.FormValue("amount")
	orderId := c.FormValue("order_id")
	t := c.Cookies("token")

	atomic, err := toAtomic(amount)
	if err != nil {
		return serveErrorPage(c, http.StatusBadRequest, err.Error())
	}

	j := &models.PostPaymentRequest{
		Amount:  atomic,
		OrderId: orderId,
	}
	// Create payment, then serve its page.
	resp, err := config.Api.PostMerchantPayment(t, j)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/merchant/request?created=true&id=" + resp.PaymentId)
}

func GetAdminPayments(c *fiber.Ctx) error {
	token := c.Cookies("token")
	resp, err := config.Api.GetAdminPayment(token)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	for i, p := range resp {
		resp[i].AmountFloat = walletrpc.XMRToDecimal(p.Amount)
		resp[i].FeeFloat = walletrpc.XMRToDecimal(p.Fee)
	}
	return c.Render("admin-payments", fiber.Map{
		"PageTitle": "Payments",
		"Payments":  resp,
	}, "layouts/admin-panel")
}
