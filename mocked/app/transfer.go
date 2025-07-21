package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Transfer struct {
	Amount int64 `json:"amount" xml:"amount" form:"amount" query:"amount"`
	Name   string	`json:"name" xml:"name" form:"name" query:"name"`
	TaxId  string  `json:"tax_id" xml:"tax_id" form:"tax_id" query:"tax_id"`
}

func MakeTransfer(c echo.Context) error {
	i := Transfer{
		Name:  "Renarin",
		Amount: 40000,
		TaxId: "155.555.555-47",
	}
	return c.JSON(http.StatusCreated, i)
}
 