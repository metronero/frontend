package server

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"

	"gitlab.com/metronero/frontend/app/controllers"
	"gitlab.com/metronero/frontend/utils/config"
)

func (s *Server) RegisterRoutes() {
	s.Static("/assets", "./public")

	s.Get("/login", controllers.GetLogin)
	s.Post("/login", controllers.PostLogin)
	s.Get("/register", controllers.GetRegister)
	s.Post("/register", controllers.PostRegister)
	s.Get("/logout", controllers.GetLogout)
	s.Get("/", controllers.GetLogin)
	s.Get("/p/:id", controllers.GetPaymentPage)

	merchant := s.Group("/merchant")
	merchant.Use(jwtware.New(jwtware.Config{
		SigningKey:  []byte(config.JwtSecret),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login?expired=true")
		},
	}))
	merchant.Use(MerchantJwtMiddleware)
	merchant.Get("/account", controllers.MerchantAccount)
	merchant.Get("/dashboard", controllers.MerchantDashboard)
	merchant.Get("/payments", controllers.MerchantPayments)
	merchant.Get("/withdrawals", controllers.MerchantWithdrawals)
	merchant.Get("/request", controllers.MerchantRequestPaymentForm)
	merchant.Post("/request", controllers.MerchantRequestPayment)
	merchant.Get("/template", controllers.MerchantTemplate)
	merchant.Get("/template/preview", controllers.MerchantTemplatePreview)
	merchant.Post("/template", controllers.MerchantTemplateUpload)
	merchant.Post("/template/reset", controllers.MerchantTemplateReset)

	admin := s.Group("/admin")
	admin.Use(jwtware.New(jwtware.Config{
		SigningKey:  []byte(config.JwtSecret),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login?expired=true")
		},
	}))
	admin.Use(AdminJwtMiddleware)
	admin.Get("/dashboard", controllers.AdminDashboard)
	admin.Get("/instance", controllers.AdminInstance)
	admin.Post("/instance", controllers.AdminEditInstance)
	admin.Get("/withdrawals", controllers.AdminWithdrawals)
	admin.Get("/payments", controllers.AdminPayments)
	admin.Get("/merchants", controllers.AdminMerchants)
	admin.Get("/merchants/:id", controllers.AdminGetMerchant)
	admin.Post("/merchants/:id", controllers.AdminEditMerchant)
	admin.Post("/merchants/delete/:id", controllers.AdminDeleteMerchant)
}
