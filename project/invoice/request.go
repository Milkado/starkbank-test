package invoice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"test/starkbank/helpers"
	"test/starkbank/project/queue"
	"time"
)

var logFile = "../logs/create_invoice.txt"

func CreateInvoice() {
	// For the example, we use 3 seconds. For your use case, change this to 3 * time.Hour.
	d := 5 * time.Second

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	var mu sync.Mutex
	var wg sync.WaitGroup
	q := &queue.QueueRequest{}

	// Create a context that will be canceled after 24 minutes.
	ctx, cancel := context.WithTimeout(context.Background(), 16*time.Second)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping the periodic task due to timeout.")
				return
			case t := <-ticker.C:
				mu.Lock()
				q = queueInvoices()
				for q.LengthRequest() > 0 {
					resp, code, err := requestQueued(*q.DequeueRequest())
					if code != http.StatusCreated {
						//log error to file
						message := fmt.Errorf("error creating invoice, returned status code: %v with message: %v. error trace: %v", strconv.Itoa(code), resp, err.Error())
						helpers.LogError(logFile, message.Error())
						continue
					}
					fmt.Println(resp)
				}
				mu.Unlock()
				fmt.Printf("Tick at %v, running the task... \n", t.Format("15:04:05"))
			}
		}
	}()

	fmt.Printf("Periodic task scheduled every %s. Will stop after 24 minutes.\n", d)

	// Wait for the context to be canceled, which happens after 24 minutes.
	<-ctx.Done()

	// Wait for the goroutine to finish.
	wg.Wait()

	fmt.Println("Invoices queued.")
	mu.Lock()
	for q.LengthRequest() > 0 {
		fmt.Println(q.DequeueRequest())
	}
	mu.Unlock()
	fmt.Println("Invoices Requested")
}

func queueInvoices() *queue.QueueRequest {
	q := &queue.QueueRequest{}

	for i := 1; i <= 12; i++ {
		q.EnqueueRequest(queue.Invoice{
			Name:   "Renarin Kholin" + strconv.Itoa(i),
			Amount: 4000 + 1,
			TaxId:  "15564426644",
		})
	}

	return q
}

func requestQueued(invoice queue.Invoice) (string, int, error) {
	invoiceJson, err := json.Marshal(invoice)
	if err != nil {
		return "", 500, err
	}

	encodedJson := bytes.NewBuffer(invoiceJson)

	createUrl := helpers.Env("STARK_API") + "/invoice"

	resp, err := http.Post(createUrl, "application/json", encodedJson)
	if err != nil {
		// If http.Post fails, resp is likely nil. Return 0 for status code.
		return "", resp.StatusCode, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	stringBody := string(respBody)
	if err != nil {
		// Return the string body even if there was a read error.
		return stringBody, resp.StatusCode, err
	}

	if resp.StatusCode != http.StatusCreated {
		errMessage := fmt.Errorf("error creating invoice")
		return stringBody, resp.StatusCode, errMessage
	}

	return stringBody, resp.StatusCode, nil
}
