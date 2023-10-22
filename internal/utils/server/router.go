package server

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"

	"gitlab.com/metronero/frontend/internal/app/controllers"
	"gitlab.com/metronero/frontend/internal/utils/config"
)

func (s *Server) RegisterRoutes() {
	s.Static("/assets", "./public")

	s.Get("/login", controllers.GetLogin)
	s.Post("/login", controllers.PostLogin)
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
	merchant.Get("/account", controllers.GetMerchantAccount)
	merchant.Get("/dashboard", controllers.GetMerchantDashboard)
	merchant.Get("/payments", controllers.GetMerchantPayments)
	merchant.Get("/withdrawals", controllers.GetMerchantWithdrawals)
	merchant.Get("/request", controllers.GetMerchantRequest)
	merchant.Post("/request", controllers.PostMerchantRequest)
	merchant.Get("/template", controllers.GetMerchantTemplate)
	merchant.Get("/template/preview", controllers.GetMerchantTemplatePreview)
	merchant.Post("/template", controllers.PostMerchantTemplate)
	merchant.Post("/template/reset", controllers.PostMerchantTemplateReset)

	admin := s.Group("/admin")
	admin.Use(jwtware.New(jwtware.Config{
		SigningKey:  []byte(config.JwtSecret),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login?expired=true")
		},
	}))
	admin.Use(AdminJwtMiddleware)
	admin.Get("/dashboard", controllers.GetAdminDashboard)
	admin.Get("/withdrawals", controllers.GetAdminWithdrawals)
	admin.Get("/payments", controllers.GetAdminPayments)
	admin.Get("/merchants", controllers.AdminMerchants)
	admin.Post("/merchants", controllers.AdminCreateMerchant)
	admin.Get("/merchants/:id", controllers.AdminGetMerchant)
	admin.Post("/merchants/:id", controllers.AdminEditMerchant)
	admin.Post("/merchants/delete/:id", controllers.AdminDeleteMerchant)
}
