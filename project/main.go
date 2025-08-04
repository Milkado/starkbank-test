package main

import (
	"context"
	"test/starkbank/config"
	"test/starkbank/project/queue"
	"test/starkbank/project/requests"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var logFile = "../logs/project_error.txt"

func main() {
	ctx := context.Background()

	cfg := config.ConfigAWS(ctx)

	sqsClient := sqs.NewFromConfig(cfg)

	newClient := queue.SqsAction(sqsClient)
	queueName := "invoices.fifo"
	var queueUrl string

	queueUrl = newClient.GetQueue(ctx, queueName)
	if queueUrl == "" {
		queueUrl = newClient.CreateSqsQueue(ctx, queueName, true)
	}

	requests.CreateInvoice(queueUrl, newClient)

}
