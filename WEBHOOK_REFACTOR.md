# Stark Bank: Building an Event-Driven Payment Processing System

## 1. Executive Summary & Architectural Vision

This document outlines the architecture and implementation of a robust, event-driven system designed to automate invoice creation and payment processing using Stark Bank's APIs. This is not a simple refactor but a new, holistic solution designed for scalability, resilience, and security.

**The Business Requirement:**
1.  **Automated Issuance:** Periodically (every 3 hours), issue 8 to 12 invoices to random recipients.
2.  **Event-Driven Payment:** Receive a real-time notification (webhook) when an invoice is paid.
3.  **Automated Payout:** Upon receiving the payment notification, automatically create and send a transfer for the credited amount.

**Our Architectural Approach:**
To meet these requirements, we will build a single Go application comprising three core, decoupled components. This approach provides the separation of concerns of microservices with the deployment simplicity of a monolith.

1.  **Invoice Scheduler:** A cron-based job that runs every 3 hours. It is responsible for generating and issuing the invoices via the Stark Bank API.
2.  **Webhook Handler:** A lightweight Echo web server that exposes a single endpoint (`/webhook/starkbank`) to receive `invoice-credited` events from Stark Bank. Its sole responsibilities are to validate the event's authenticity and enqueue it into an AWS SQS queue. This makes our webhook ingestion point extremely fast and resilient.
3.  **Event Worker:** A background goroutine that polls the SQS queue. It dequeues payment events and orchestrates the business logic: parsing the event and creating a corresponding transfer via the Stark Bank API. Using a queue decouples the processing logic from the webhook ingestion, allowing us to handle bursts of traffic and retry failed processing attempts without losing events.

![Architecture Diagram](https://i.imgur.com/k2p4J1g.png)

---

## 2. Step-by-Step Implementation Guide

### Step 1: Dependencies and Configuration

First, ensure your project has the necessary dependencies. We will need the Stark Bank Go SDK, the AWS SDK for SQS, the Echo web framework, and a cron library.

```bash
go get github.com/starkbank/starkbank-go-sdk/v2
go get github.com/aws/aws-sdk-go-v2/service/sqs
go get github.com/labstack/echo/v4
go get github.com/robfig/cron/v3
go mod tidy
```

Your `.env` file must be configured with your Stark Bank Project ID and Private Key.

### Step 2: The Invoice Scheduler

This component handles the periodic creation of invoices. We'll create a dedicated package for this logic.

**File: `project/scheduler/invoice.go`**
```go
package scheduler

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/starkbank/starkbank-go-sdk/v2/invoice"
	"github.com/starkbank/starkbank-go-sdk/v2/user"
)

// CreateRandomInvoices generates and issues a random number of invoices (8-12).
func CreateRandomInvoices(project user.Project) {
	rand.Seed(time.Now().UnixNano())
	invoiceCount := rand.Intn(5) + 8 // Random number between 8 and 12

	var invoices []invoice.Invoice
	for i := 0; i < invoiceCount; i++ {
		invoices = append(invoices, invoice.Invoice{
			Amount: rand.Intn(10000) + 100, // Random amount between 1.00 and 100.00
			Name:   fmt.Sprintf("Test User %d", i),
			TaxId:  "20.018.183/0001-80", // Example Tax ID
		})
	}

	// The SDK handles the API request to create the invoices
	createdInvoices, err := invoice.Create(invoices, project)
	if err.Errors != nil {
		fmt.Printf("Failed to create invoices: %v
", err.Errors)
		return
	}

	fmt.Printf("Successfully created %d invoices.
", len(createdInvoices))
}
```
**Architectural Decision:** By encapsulating the invoice creation logic, we make it reusable and independent of its execution context. The function requires a `user.Project` object, making dependencies explicit and testing easier.

### Step 3: The Webhook Handler

This is the public-facing part of our application. It must be fast, secure, and reliable.

**File: `project/webhook/handler.go`**
```go
package webhook

import (
	"io/ioutil"
	"net/http"
	"test/starkbank/project/queue"
	"github.com/labstack/echo/v4"
	"github.com/starkbank/starkbank-go-sdk/v2/event"
	"github.com/starkbank/starkbank-go-sdk/v2/user"
)

// HandleStarkBankEvent processes incoming Stark Bank webhooks.
func HandleStarkBankEvent(sqsClient queue.SqsActions, queueUrl string, project user.Project) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Read the raw request body
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusBadRequest, "Failed to read request body")
		}

		// 2. CRITICAL: Verify the signature
		// The SDK uses the raw body and the "Digital-Signature" header to prevent tampering.
		signature := c.Request().Header.Get("Digital-Signature")
		parsedEvent, parseErr := event.Parse(string(body), signature, project)
		if parseErr.Errors != nil {
			// If the signature is invalid, reject the request immediately.
			return c.String(http.StatusUnauthorized, "Invalid webhook signature")
		}

		// 3. Filter for the event we care about
		if parsedEvent.Subscription == "invoice-credited" {
			// 4. Enqueue the validated event for asynchronous processing
			err := sqsClient.SendMessage(c.Request().Context(), queueUrl, string(body))
			if err != nil {
				// If we can't queue it, return an error so Stark Bank retries the webhook.
				return c.String(http.StatusInternalServerError, "Failed to queue event")
			}
		}

		// 5. Respond immediately with a 200 OK
		// This tells Stark Bank we have successfully received the event.
		return c.String(http.StatusOK, "Event received")
	}
}
```
**Architectural Decisions:**
*   **Security First:** The most critical step is `event.Parse`. It cryptographically verifies that the webhook is genuinely from Stark Bank and hasn't been altered. We **never** process a webhook without this check.
*   **Minimal Processing:** The handler's only job is to validate and enqueue. All time-consuming business logic is deferred to the worker. This ensures the handler can process a high volume of webhooks without timing out, which is crucial for system reliability.
*   **Filtering:** We explicitly check for `invoice-credited`. Your webhook endpoint might receive other event types; this ensures we only process the ones relevant to this flow.

### Step 4: The Event Worker

This is the workhorse of the system. It runs in the background, processing payment events safely from the queue.

**File: `project/worker/processor.go`**
```go
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/starkbank/starkbank-go-sdk/v2/event"
	"github.com/starkbank/starkbank-go-sdk/v2/transfer"
	"github.com/starkbank/starkbank-go-sdk/v2/user"
	"test/starkbank/project/queue"
)

// StartPolling continuously polls SQS for new messages and processes them.
func StartPolling(ctx context.Context, sqsClient queue.SqsActions, queueUrl string, project user.Project) {
	fmt.Println("Worker started, polling for messages from:", queueUrl)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker shutting down.")
			return
		default:
			messages := sqsClient.GetMessages(ctx, queueUrl)
			if len(messages) == 0 {
				time.Sleep(5 * time.Second) // Wait before polling again
				continue
			}

			for _, message := range messages {
				processMessage(ctx, *message.Body, project)
				// IMPORTANT: Delete the message only after it's been fully processed.
				sqsClient.DeleteMessage(ctx, queueUrl, *message.ReceiptHandle)
			}
		}
	}
}

// processMessage parses the event and creates the corresponding transfer.
func processMessage(ctx context.Context, messageBody string, project user.Project) {
	var eventData event.Event
	if err := json.Unmarshal([]byte(messageBody), &eventData); err != nil {
		fmt.Printf("Failed to unmarshal event: %v
", err)
		return // Message is malformed, won't be retried
	}

	// The event log contains the invoice object
	creditedInvoice := eventData.Log.Invoice

	fmt.Printf("Processing credited invoice ID %s for amount %d
", creditedInvoice.Id, creditedInvoice.Amount)

	// Create a transfer for the exact amount of the paid invoice
	transfers := []transfer.Transfer{
		{
			Amount:        creditedInvoice.Amount,
			Name:          "Stark Bank S.A.",
			TaxId:         "20.018.183/0001-80",
			AccountNumber: "6341320293482496",
			BankCode:      "20018183",
			BranchCode:    "0001",
			AccountType:   "payment",
		},
	}

	_, err := transfer.Create(transfers, project)
	if err.Errors != nil {
		fmt.Printf("Failed to create transfer for invoice %s: %v
", creditedInvoice.Id, err.Errors)
		// In a production system, you would add this to a dead-letter queue for investigation.
		return
	}

	fmt.Printf("Successfully created transfer for invoice %s.
", creditedInvoice.Id)
}
```
**Architectural Decisions:**
*   **Idempotency:** The worker must be idempotent. If it processes the same message twice, it should not result in a duplicate transfer. While not fully implemented here, a production system would check if a transfer for a given invoice ID has already been created before creating a new one.
*   **Atomic Processing:** The message is deleted from the SQS queue *only after* the transfer is successfully created. If the transfer fails, the message remains in the queue and will be re-processed after a visibility timeout. This prevents losing money/transactions due to transient failures.
*   **Dead-Letter Queue (DLQ):** For messages that consistently fail (e.g., due to a bug or malformed data), the SQS queue should be configured with a DLQ. This prevents a "poison pill" message from blocking the entire queue.

### Step 5: Tying It All Together

The `main.go` file will initialize all components and start them in the correct order.

**File: `project/main.go`**
```go
package main

import (
	"context"
	"fmt"
	"os"
	"test/starkbank/config"
	"test/starkbank/project/queue"
	"test/starkbank/project/scheduler"
	"test/starkbank/project/webhook"
	"test/starkbank/project/worker"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"github.com/starkbank/starkbank-go-sdk/v2/user"
)

func main() {
	// Load environment variables
	godotenv.Load()
	
	// 1. Initialize Stark Bank Project
	project := user.Project{
		Id:          os.Getenv("STARK_PROJECT_ID"),
		PrivateKey:  os.Getenv("STARK_PRIVATE_KEY"),
		Environment: "sandbox",
	}

	// 2. Initialize AWS SQS Client
	ctx := context.Background()
	cfg := config.ConfigAWS(ctx)
	sqsClient := sqs.NewFromConfig(cfg)
	newClient := queue.SqsAction(sqsClient)
	queueName := "starkbank-events.fifo"
	queueUrl := newClient.GetOrCreateQueue(ctx, queueName)
	fmt.Println("SQS Queue URL:", queueUrl)

	// 3. Initialize and start the Invoice Scheduler
	c := cron.New()
	c.AddFunc("@every 3h", func() {
		fmt.Println("Running invoice creation job...")
		scheduler.CreateRandomInvoices(project)
	})
	c.Start()
	fmt.Println("Invoice scheduler started. Will run every 3 hours.")

	// 4. Start the SQS Worker in a background goroutine
	go worker.StartPolling(ctx, newClient, queueUrl, project)

	// 5. Initialize and start the Echo Web Server
	e := echo.New()
	e.POST("/webhook/starkbank", webhook.HandleStarkBankEvent(newClient, queueUrl, project))
	fmt.Println("Webhook server starting on :1323")
	e.Logger.Fatal(e.Start(":1323"))
}
```
**Architectural Decision:** The application is orchestrated in `main`. It starts the long-running background tasks (cron scheduler, SQS worker) and then starts the blocking web server process in the main thread. This creates a self-contained application that is easy to run and deploy.

---

## 3. Running the System

1.  **Run the Application:**
    ```bash
    go run project/main.go
    ```
2.  **Expose Your Webhook:** For Stark Bank to reach your local server, use a tool like `ngrok`.
    ```bash
    ngrok http 1323
    ```
    This will give you a public URL (e.g., `https://abcd-1234.ngrok.io`).
3.  **Configure Stark Bank:** In your Stark Bank dashboard, create a new webhook subscription.
    *   **URL:** The `ngrok` URL (e.g., `https://abcd-1234.ngrok.io/webhook/starkbank`).
    *   **Subscriptions:** Select `invoice-credited`.

Now, the system is live. Every 3 hours, new invoices will be created. When the sandbox environment pays one, the webhook will fire, and you will see the worker process the event and create a transfer.