package middleware

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/pranish23/mini-aspire-app/users"
)

func CheckAdminAuthorization() func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Users: users.Admins,
	})
}

func CheckCustomerAuthorization() func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Users:           users.Customers,
		ContextUsername: "username",
	})

}
