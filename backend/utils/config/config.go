package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var JWTSecret string
var Logger *slog.Logger
var UserService *services.UserService
var DocumentService *services.DocumentService
var S3Client *s3.Client

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := 5432
	dbUser := "postgres"
	dbName := "postgres"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	userRepository := repositories.NewUserRepository(db)
	UserService = services.NewUserService(userRepository, Logger)

	S3Client, err = newS3Client()
	if err != nil {
		panic(err)
	}

	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("S3_BUCKET_NAME not set")
	}

	documentRepository := repositories.NewDocumentRepository(db)
	DocumentService = services.NewDocumentService(documentRepository, userRepository, S3Client, bucketName, Logger)
}

func newS3Client() (*s3.Client, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("eu-west-1"),
	)
	if err != nil {
		return nil, err
	}

	stsClient := sts.NewFromConfig(cfg)

	roleProvider := stscreds.NewAssumeRoleProvider(
		stsClient,
		"arn:aws:iam::659775407830:role/s3-doclane-role",
	)

	assumedCfg := cfg
	assumedCfg.Credentials = aws.NewCredentialsCache(roleProvider)

	s3Client := s3.NewFromConfig(assumedCfg)

	return s3Client, nil
}
