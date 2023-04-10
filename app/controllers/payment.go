package controllers

import (
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/metronero/backend/app/models"

	"gitlab.com/metronero/frontend/utils/config"
	"gitlab.com/metronero/frontend/utils/token"
)

func GetPaymentPage(c *fiber.Ctx) error {
	resp, err := config.Api.GetPaymentPage(c.Params("id"))
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	c.Set("Refresh", "10")
	return c.Send(resp)
}

func MerchantRequestPaymentForm(c *fiber.Ctx) error {
	id := c.Query("id", "")
	created := c.Query("created", "false")

	var (
		reqUrl string
		err error
	)
	if id != "" {
		reqUrl, err = url.JoinPath(config.Hostname, "p", id)
		if err != nil {
			return serveErrorPage(c, http.StatusInternalServerError, err.Error())
		}
	}
	return c.Render("merchant-request-payment", fiber.Map{
		"Username":  token.GetUsername(c),
		"PageTitle": "Request Payment",
		"PaymentReqUrl": reqUrl,
		"Created": created,
	}, "layouts/merchant-panel")
}

func MerchantRequestPayment(c *fiber.Ctx) error {
	amount := c.FormValue("amount")
	orderId := c.FormValue("order_id")
	t := c.Cookies("token")

	atomic, err := toAtomic(amount)
	if err != nil {
		return serveErrorPage(c, http.StatusBadRequest, err.Error())
	}

	j := &models.PostPaymentRequest{
		Amount: atomic,
		OrderId: orderId,
	}
	// Create payment, then serve its page.
	resp, err := config.Api.CreatePaymentRequest(t, j)
	if err != nil {
		return serveErrorPage(c, http.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/merchant/request?created=true&id="+resp.PaymentId)
}
