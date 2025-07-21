package model

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type (
	Invoice struct {
		Amount     float64
		TaxId      string
		Due        time.Time
		Expiration int64
		Fine       float64
		Interest   float64
		Fee        float64
		Status     string
	}

	InviceResp struct {
		ID         int64
		Amount     float64
		TaxId      string
		Due        time.Time
		Expiration int64
		Fine       float64
		Interest   float64
		Fee        float64
		Status     string
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	InvoiceRequest struct {
		Amount float64 `json:"amount" xml:"amount" form:"amount" query:"amount"`
		Name   string  `json:"name" xml:"name" form:"name" query:"name"`
		TaxId  string  `json:"tax_id" xml:"tax_id" form:"tax_id" query:"tax_id"`
	}
)

// query example
func InvoiceById(id int64, db *sql.DB) (InviceResp, error) {
	row := db.QueryRow("SELECT * FROM invoice WHERE id = ?", id)

	invoiceResp := InviceResp{}
	if err := row.Scan(&invoiceResp.ID, &invoiceResp.Amount, &invoiceResp.TaxId, &invoiceResp.Due, &invoiceResp.Expiration, &invoiceResp.Fine, &invoiceResp.Interest, &invoiceResp.Fee, &invoiceResp.Status, &invoiceResp.CreatedAt, &invoiceResp.UpdatedAt); err != nil {
		return InviceResp{}, fmt.Errorf("no invoice with this Id %q: %v", id, err)
	}
	invoiceResp.Status = getStatus(invoiceResp.Status)

	return invoiceResp, nil
}

func StoreInvoice(request InvoiceRequest, db *sql.DB) (InviceResp, error) {
	invoice := Invoice{
		Amount:     request.Amount,
		TaxId:      request.TaxId,
		Due:        time.Now().AddDate(0, 0, 4),
		Expiration: 4,
		Fine:       0,
		Interest:   2,
		Fee:        3.4,
		Status:     randomStatus(),
	}

	result, err := db.Exec("INSERT INTO invoice (amount, tax_id, due, expiration, fine, interest, fee, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", invoice.Amount, invoice.TaxId, invoice.Due, invoice.Expiration, invoice.Fine, invoice.Interest, invoice.Fee, invoice.Status)
	if err != nil {
		return InviceResp{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return InviceResp{}, err
	}

	row := db.QueryRow("SELECT * FROM invoice WHERE id = ?", id)

	invoiceResp := InviceResp{}
	if err := row.Scan(&invoiceResp.ID, &invoiceResp.Amount, &invoiceResp.TaxId, &invoiceResp.Due, &invoiceResp.Expiration, &invoiceResp.Fine, &invoiceResp.Interest, &invoiceResp.Fee, &invoiceResp.Status, &invoiceResp.CreatedAt, &invoiceResp.UpdatedAt); err != nil {
		return InviceResp{}, fmt.Errorf("error getting Invoice with this Id %q: %v", id, err)
	}
	invoiceResp.Status = getStatus(invoiceResp.Status)
	return invoiceResp, nil
}

func randomStatus() string {
	status := make(map[int]string)
	status[1] = "P"
	status[2] = "C"

	min, max := 1, 2

	return status[min+rand.Intn(max-min)]
}

func getStatus(char string) string {
	status := make(map[string]string)
	status["P"] = "PAYED"
	status["C"] = "CREATED"

	return status[char]
}
