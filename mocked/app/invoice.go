package app

import (
	"fmt"
	"net/http"
	"strconv"
	"test/starkbank/helpers"
	"test/starkbank/mocked/app/model"
	"test/starkbank/mocked/db"

	"github.com/labstack/echo/v4"
)

var dbConn = db.DbConn{
	User:   helpers.Env("DB_USER"),
	Pass:   helpers.Env("DB_PASSWORD"),
	Addr:   helpers.Env("DB_HOST") + ":" + helpers.Env("DB_PORT"),
	DbName: helpers.Env("DB_NAME"),
}

func CreateInvoice(c echo.Context) error {
	conn, err := db.Connect(dbConn)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	i := new(model.InvoiceRequest)
	if bindErr := c.Bind(i); bindErr != nil {
		return c.JSON(http.StatusBadRequest, bindErr.Error())
	}

	if i.Name == "Renarin Kholin12" || i.Name == "Renarin Kholin2" {
		ficErr := fmt.Errorf("error for testing")
		return c.JSON(http.StatusBadRequest, ficErr.Error())
	}

	resp, err := model.StoreInvoice(*i, conn)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, resp)
}

func ConsultInvoice(c echo.Context) error {
	conn, err := db.Connect(dbConn)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("error converting invoice id to int: ", err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	resp, err := model.InvoiceById(int64(id), conn)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
