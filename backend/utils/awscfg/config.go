package awscfg

import (
	"context"
	"log"

	"github.com/Robert076/doclane/backend/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func InitAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(utils.RequireEnv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal("Failed to load AWS config:", err)
	}
	return cfg
}
