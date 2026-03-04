package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

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
var InvitationCodeService *services.InvitationCodeService
var ExpectedDocumentService *services.ExpectedDocumentService
var DocumentRequestTemplateService *services.DocumentRequestTemplateService
var S3Client *s3.Client

func init() {
	// Load .env file if present — ignored in Lambda since env vars are injected directly
	godotenv.Load("../../.env")

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db := initDB()
	JWTSecret = requireEnv("JWT_SECRET")

	// Initialize repositories
	userRepository := repositories.NewUserRepository(db)
	documentRepository := repositories.NewDocumentRepository(db)
	invitationRepository := repositories.NewInvitationCodeRepository(db)
	expectedDocumentRepo := repositories.NewExpectedDocRepository(db)
	documentRequestTemplateRepo := repositories.NewDocumentRequestTemplateRepository(db)
	expectedDocumentTemplateRepo := repositories.NewExpectedDocumentTemplateRepository(db)
	txManager := repositories.NewTxManager(db)

	// Initialize S3
	var err error
	S3Client, err = newS3Client()
	if err != nil {
		log.Fatal("Failed to initialize S3 client:", err)
	}

	// Initialize services
	fileStorage := services.NewFileStorageService(S3Client, requireEnv("S3_BUCKET_NAME"), Logger)
	UserService = services.NewUserService(userRepository, Logger)
	DocumentService = services.NewDocumentService(documentRepository, userRepository, expectedDocumentRepo, txManager, Logger, fileStorage)
	InvitationCodeService = services.NewInvitationCodeService(invitationRepository, userRepository, Logger)
	ExpectedDocumentService = services.NewExpectedDocumentService(expectedDocumentRepo, Logger)
	DocumentRequestTemplateService = services.NewDocumentRequestTemplateService(
		documentRequestTemplateRepo,
		expectedDocumentTemplateRepo,
		expectedDocumentRepo,
		documentRepository,
		userRepository,
		txManager,
		fileStorage,
		Logger,
	)
}

func initDB() *sql.DB {
	host := requireEnv("DB_HOST")
	user := requireEnv("DB_USER")
	name := requireEnv("DB_NAME")
	password := requireEnv("DB_PASSWORD")

	port := 5432
	if p := os.Getenv("DB_PORT"); p != "" {
		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			log.Fatal("DB_PORT must be a number")
		}
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, name, password,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	return db
}

func newS3Client() (*s3.Client, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-1"))
	if err != nil {
		return nil, err
	}

	s3IamRole := requireEnv("AWS_ROLE_S3")

	roleProvider := stscreds.NewAssumeRoleProvider(sts.NewFromConfig(cfg), s3IamRole)

	assumedCfg := cfg
	assumedCfg.Credentials = aws.NewCredentialsCache(roleProvider)

	return s3.NewFromConfig(assumedCfg), nil
}

// requireEnv gets an env var and fatals if it's not set
func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s env var is not set", key)
	}
	return val
}
