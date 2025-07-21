package routes

import (
	"test/starkbank/mocked/app"

	"github.com/labstack/echo/v4"
)

func Api() *echo.Echo {
	e := echo.New()

	e.POST("/invoice", app.CreateInvoice)

	return e
}
