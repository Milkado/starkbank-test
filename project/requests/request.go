package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"test/starkbank/helpers"
	"test/starkbank/project/queue"
	"time"

	"github.com/gosimple/slug"
)

var logFile = "../logs/create_invoice.txt"

type (
	Invoice struct {
		Amount float64 `json:"amount" xml:"amount" form:"amount" query:"amount"`
		Name   string  `json:"name" xml:"name" form:"name" query:"name"`
		TaxId  string  `json:"tax_id" xml:"tax_id" form:"tax_id" query:"tax_id"`
	}
)

func CreateInvoice(queueUrl string, sqsClient queue.SqsActions) {
	// For the example, we use 3 seconds. For your use case, change this to 3 * time.Hour.
	d := 3 * time.Minute

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create a context that will be canceled after 24 minutes.
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Minute)
	defer cancel()

	i := 0

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
				i++
				queueInvoices(ctx, i, sqsClient, queueUrl)
				fmt.Printf("Sleeping for 30s at %v \n", t.Format("15:04:05"))
				time.Sleep(30*time.Second)
				messages := sqsClient.GetMessages(ctx, queueUrl)
				var invoice Invoice
				for _, message := range messages {
					err := json.Unmarshal([]byte(*message.Body), &invoice)
					if err != nil {
						helpers.LogError(logFile, err.Error())
						continue
					}

					requestCreation(invoice)
					sqsClient.DeleteMessage(ctx, queueUrl, *message.ReceiptHandle)
				}
				mu.Unlock()
				fmt.Printf("Tick at %v, running the task... \n", t.Format("15:04:05"))
			}
		}
	}()

	fmt.Printf("Periodic task scheduled every %m. Will stop after 24 minutes.\n", d)

	// Wait for the context to be canceled, which happens after 24 minutes.
	<-ctx.Done()
	// Wait for the goroutine to finish.
	wg.Wait()
	fmt.Println("Invoices queued.")

	fmt.Println("Invoices Requested")
}

func queueInvoices(ctx context.Context, requestId int, sqsClient queue.SqsActions, queueUrl string) {
	for y := 1; y <= 8; y++ {
		newInvoice := Invoice{
			Amount: 4000.10 + float64(requestId),
			Name:   "Kaladin Stormblessed",
			TaxId:  "11111111111",
		}
		message, err := json.Marshal(newInvoice)
		if err != nil {
			helpers.LogError(logFile, err.Error())
		}
		g := new(string)
		duplicationId := new(string)
		*duplicationId = slug.Make(newInvoice.Name + strconv.Itoa(requestId) + strconv.Itoa(y))
		*g = strconv.Itoa(requestId)

		sqsClient.SendMessage(ctx, queueUrl, message, g, duplicationId)
	}
}

func requestCreation(invoice Invoice) {
	fmt.Println(invoice)
}
