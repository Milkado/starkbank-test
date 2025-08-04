package config

import (
	"context"
	"os"
	"test/starkbank/helpers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var logFile = "/logs/aws_config.txt"

func ConfigAWS(ctx context.Context) aws.Config {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"), config.WithSharedConfigProfile("AdminDev"))
	if err != nil {
		helpers.LogError(logFile, err.Error())
		os.Exit(1)
	}

	return cfg
}