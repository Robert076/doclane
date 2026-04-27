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
	"github.com/Robert076/doclane/backend/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Logger *slog.Logger
var UserService *services.UserService
var RequestService *services.RequestService
var DepartmentService *services.DepartmentService
var InvitationCodeService *services.InvitationCodeService
var ExpectedDocumentService *services.ExpectedDocumentService
var RequestTemplateService *services.RequestTemplateService
var RequestCommentService *services.RequestCommentService
var StatsService *services.StatsService
var TagService *services.TagService
var S3Client *s3.Client

func init() {
	godotenv.Load("../../.env")

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db := initDB()

	// Repositories
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

	// S3
	var err error
	S3Client, err = newS3Client()
	if err != nil {
		log.Fatal("Failed to initialize S3 client:", err)
	}

	// Services
	fileStorage := services.NewFileStorageService(S3Client, utils.RequireEnv("S3_BUCKET_NAME"), Logger)

	UserService = services.NewUserService(userRepo, Logger)
	RequestService = services.NewRequestService(
		requestRepo,
		userRepo,
		requestTemplateRepo,
		expectedDocumentRepo,
		expectedDocumentTemplateRepo,
		txManager,
		Logger,
		fileStorage,
	)
	DepartmentService = services.NewDepartmentService(departmentRepo, Logger)
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

func initDB() *sql.DB {
	host := utils.RequireEnv("DB_HOST")
	user := utils.RequireEnv("DB_USER")
	name := utils.RequireEnv("DB_NAME")
	password := utils.RequireEnv("DB_PASSWORD")

	port := 5432
	if p := os.Getenv("DB_PORT"); p != "" {
		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			log.Fatal("DB_PORT must be a number")
		}
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable timezone=UTC",
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
	return s3.NewFromConfig(cfg), nil
}
