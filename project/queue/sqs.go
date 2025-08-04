package queue

import (
	"context"
	"log"
	"os"
	"test/starkbank/helpers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var logFile = "../logs/sqs_errors.txt"

type SqsActions struct {
	SqsClient *sqs.Client
}


func SqsAction(sqsCliet *sqs.Client) SqsActions {
	return SqsActions{
		SqsClient: sqsCliet,
	}
}

func (actor SqsActions) CreateSqsQueue(ctx context.Context, queueName string, isFifo bool) string {
	queueAttributes := map[string]string{}

	if isFifo {
		queueAttributes["FifoQueue"] = "true"
	}

	queue, err := actor.SqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName:  aws.String(queueName),
		Attributes: queueAttributes,
	})
	if err != nil {
		helpers.LogError(logFile, err.Error())
		os.Exit(1)
	}

	return *queue.QueueUrl
}

func (actor SqsActions) GetQueue(ctx context.Context, queueName string) string {
	queue, err := actor.SqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
		QueueOwnerAWSAccountId: nil,
	})
	if err != nil {
		return ""
	}

	return *queue.QueueUrl
}

func (actor SqsActions) SendMessage(ctx context.Context, queueUrl string, message []byte, group *string, dupId *string) {
	res, err := actor.SqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:             aws.String(string(message)),
		QueueUrl:                &queueUrl,
		DelaySeconds:            0,
		MessageAttributes:       nil,
		MessageDeduplicationId:  dupId,
		MessageGroupId:          group,
		MessageSystemAttributes: nil,
	})
	if err != nil {
		helpers.LogError(logFile, err.Error())
		os.Exit(1)
	}

	log.Printf("the message with id %v is sent\n", *res.MessageId)
}

func (actor SqsActions) GetMessages(ctx context.Context, queueUrl string) []types.Message {
	res, err := actor.SqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueUrl),
		MaxNumberOfMessages: 8,
		WaitTimeSeconds: 20,
	})
	if err != nil {
		helpers.LogError(logFile, err.Error())
		os.Exit(1)
	}
	
	return res.Messages
}

func (actor SqsActions) PurgeQueue(ctx context.Context, queueUrl string) {
	_, err := actor.SqsClient.PurgeQueue(ctx, &sqs.PurgeQueueInput{
		QueueUrl: &queueUrl,
	})
	if err != nil {
		helpers.LogError(logFile, err.Error())
		os.Exit(1)
	}
}

func (actor SqsActions) DeleteMessage(ctx context.Context, queueUrl string, handle string) {
	_, err := actor.SqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl: &queueUrl,
		ReceiptHandle: &handle,
	})
	if err != nil {
		helpers.LogError(logFile, err.Error())
		os.Exit(1)
	}
}

