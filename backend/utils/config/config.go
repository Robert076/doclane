package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/Robert076/doclane/backend/events"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/services"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/awscfg"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Logger *slog.Logger
var AuditLogService *services.AuditLogService
var UserService *services.UserService
var RequestService *services.RequestService
var DepartmentService *services.DepartmentService
var InvitationCodeService *services.InvitationCodeService
var ExpectedDocumentService *services.ExpectedDocumentService
var RequestTemplateService *services.RequestTemplateService
var RequestCommentService *services.RequestCommentService
var StatsService *services.StatsService
var TagService *services.TagService
var TextractService *services.TextractService
var BedrockService *services.BedrockService
var PollyService *services.PollyService
var S3Client *s3.Client
var (
	AWSRegion         string
	CognitoUserPoolID string
	CognitoClientID   string
)

func init() {
	godotenv.Load("../../.env")

	// AWS - environment variables
	AWSRegion = utils.RequireEnv("AWS_REGION")
	CognitoUserPoolID = utils.RequireEnv("COGNITO_USER_POOL_ID")
	CognitoClientID = utils.RequireEnv("COGNITO_CLIENT_ID")

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	awsCfg := awscfg.InitAWSConfig()
	db := initDB(awsCfg)

	// Repositories
	auditLogRepo := repositories.NewAuditLogRepo(db)
	userRepo := repositories.NewUserRepo(db)
	requestRepo := repositories.NewRequestRepo(db)
	departmentRepo := repositories.NewDepartmentRepo(db)
	invitationRepo := repositories.NewInvitationCodeRepo(db)
	expectedDocumentRepo := repositories.NewExpectedDocRepo(db)
	requestTemplateRepo := repositories.NewRequestTemplateRepo(db)
	expectedDocumentTemplateRepo := repositories.NewExpectedDocumentTemplateRepo(db)
	requestCommentRepo := repositories.NewRequestCommentRepo(db)
	statsRepo := repositories.NewStatsRepo(db)
	tagRepo := repositories.NewTagRepo(db)
	txManager := repositories.NewTxManager(db)

	// AWS - clients
	S3Client = s3.NewFromConfig(awsCfg)
	textractClient := textract.NewFromConfig(awsCfg)
	bedrockClient := bedrockruntime.NewFromConfig(awsCfg)
	pollyClient := polly.NewFromConfig(awsCfg)

	bucket := utils.RequireEnv("S3_BUCKET_NAME")

	// Services
	fileStorage := services.NewFileStorageService(S3Client, bucket, Logger)
	TextractService = services.NewTextractService(textractClient, bucket, Logger)
	BedrockService = services.NewBedrockService(bedrockClient, Logger)
	PollyService = services.NewPollyService(pollyClient, Logger)

	// Event bus, generic, shared across services
	eventBus := events.NewEventBus(Logger)

	AuditLogService = services.NewAuditLogService(auditLogRepo, Logger)
	eventBus.Subscribe(AuditLogService)
	UserService = services.NewUserService(userRepo, requestRepo, Logger, eventBus)
	RequestService = services.NewRequestService(
		requestRepo,
		userRepo,
		requestTemplateRepo,
		expectedDocumentRepo,
		expectedDocumentTemplateRepo,
		txManager,
		Logger,
		fileStorage,
		TextractService,
		BedrockService,
		PollyService,
		eventBus,
	)
	DepartmentService = services.NewDepartmentService(departmentRepo, Logger, eventBus)
	InvitationCodeService = services.NewInvitationCodeService(invitationRepo, departmentRepo, Logger)
	ExpectedDocumentService = services.NewExpectedDocumentService(expectedDocumentRepo, requestRepo, Logger)
	RequestTemplateService = services.NewRequestTemplateService(
		requestTemplateRepo,
		expectedDocumentTemplateRepo,
		expectedDocumentRepo,
		requestRepo,
		txManager,
		fileStorage,
		Logger,
	)
	RequestCommentService = services.NewRequestCommentService(
		requestCommentRepo,
		requestRepo,
		Logger,
	)
	StatsService = services.NewStatsService(statsRepo, Logger)
	TagService = services.NewTagService(tagRepo, Logger)
}

func initDB(cfg aws.Config) *sql.DB {
	var host string
	var user string
	var name string
	var password string
	var port string
	sslMode := "disable"

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		ssmClient := ssm.NewFromConfig(cfg)

		host = utils.RequireEnv("DB_HOST")
		port = utils.RequireEnv("DB_PORT")
		name = utils.RequireEnv("DB_NAME")
		sslMode = "require"

		usernamePath := utils.RequireEnv("SSM_USERNAME_PATH")
		passwordPath := utils.RequireEnv("SSM_PASSWORD_PATH")

		userParam, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
			Name:           aws.String(usernamePath),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Fatalf("failed to get SSM username param: %v", err)
		}
		user = aws.ToString(userParam.Parameter.Value)

		passParam, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
			Name:           aws.String(passwordPath),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Fatalf("failed to get SSM password param: %v", err)
		}
		password = aws.ToString(passParam.Parameter.Value)
	} else {
		host = utils.RequireEnv("DB_HOST")
		port = utils.RequireEnv("DB_PORT")
		user = utils.RequireEnv("DB_USER")
		name = utils.RequireEnv("DB_NAME")
		password = utils.RequireEnv("DB_PASSWORD")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s timezone=UTC",
		host, port, user, name, password, sslMode,
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
