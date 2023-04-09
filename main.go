package main

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"

	"gitlab.com/metronero/frontend/app/controllers"
	"gitlab.com/metronero/frontend/utils/config"
	"gitlab.com/metronero/frontend/utils/server"
)

func main() {
	config.Load()
	app := server.Init()

	app.Static("/assets", "./public")

	app.Get("/login", controllers.GetLogin)
	app.Post("/login", controllers.PostLogin)
	app.Get("/register", controllers.GetRegister)
	app.Post("/register", controllers.PostRegister)
	app.Get("/logout", controllers.GetLogout)
	app.Get("/", controllers.GetLogin)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  []byte(config.JwtSecret),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login?expired=true")
		},
	}))

	merchant := app.Group("/merchant")
	merchant.Use(server.MerchantJwtMiddleware)
	merchant.Get("/account", controllers.MerchantAccount)
	merchant.Get("/dashboard", controllers.MerchantDashboard)
	merchant.Get("/payments", controllers.MerchantPayments)
	merchant.Get("/withdrawals", controllers.MerchantWithdrawals)
	merchant.Get("/request", controllers.MerchantRequestPayment)
	merchant.Get("/template", controllers.MerchantTemplate)
	merchant.Get("/template/preview", controllers.MerchantTemplatePreview)
	merchant.Post("/template", controllers.MerchantTemplateUpload)
	merchant.Post("/template/reset", controllers.MerchantTemplateReset)

	admin := app.Group("/admin")
	admin.Use(server.AdminJwtMiddleware)
	admin.Get("/dashboard", controllers.AdminDashboard)
	admin.Get("/instance", controllers.AdminInstance)
	admin.Post("/instance", controllers.AdminEditInstance)
	admin.Get("/withdrawals", controllers.AdminWithdrawals)
	admin.Get("/payments", controllers.AdminPayments)
	admin.Get("/merchants", controllers.AdminMerchants)
	admin.Get("/merchants/:id", controllers.AdminGetMerchant)
	admin.Post("/merchants/:id", controllers.AdminEditMerchant)
	admin.Post("/merchants/delete/:id", controllers.AdminDeleteMerchant)

	app.StartServerWithGracefulShutdown()
}
