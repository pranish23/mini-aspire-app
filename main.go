package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/pranish23/mini-aspire-app/controllers"
	"github.com/pranish23/mini-aspire-app/db"
	"github.com/pranish23/mini-aspire-app/middleware"
)

func main() {
	dbConn, err := db.Init()
	if err != nil {
		return
	}
	app := fiber.New()
	loanService := fiber.New()

	lc := controllers.NewLoanController(dbConn)

	app.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("Health check"))
		return nil
	})

	app.Mount("/loan", loanService)
	loanService.Patch("/:loanID/approve", middleware.CheckAdminAuthorization(), lc.Approve)
	loanService.Post("/:customerID/submit", middleware.CheckCustomerAuthorization(), lc.Submit)
	loanService.Get("/:customerID/view", middleware.CheckCustomerAuthorization(), lc.View)
	loanService.Patch("/:customerID/repay", middleware.CheckCustomerAuthorization(), lc.AddRepayment)
	app.Listen(":3000")
}
