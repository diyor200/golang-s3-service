package app

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang-diplom-work/config"
	v1 "golang-diplom-work/internal/controller/http/v1"
	"golang-diplom-work/internal/repo"
	"golang-diplom-work/internal/service"
	"golang-diplom-work/pkg/custom-validator"
	"golang-diplom-work/pkg/hasher"
	"golang-diplom-work/pkg/httpserver"
	"golang-diplom-work/pkg/postgres"
	"golang-diplom-work/pkg/storage"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	// Configuration
	cfg, err := config.NewConfig()
	fmt.Println(cfg, "config")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Logger
	SetLogrus(cfg.Log.Level)

	// Repositories
	log.Info("Initializing postgres...")
	//url := fmt.Sprintf("user=postgres password=2001 host=localhost port=5432 dbname=postgres sslmode=disable pool_max_conns=10")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	// set up s3
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(cfg.S3.AwsRegion),
			Credentials: credentials.NewStaticCredentials(
				cfg.S3.AwsAccessKeyId,
				cfg.S3.AwsSecretAccessKey,
				""),
		},
	})
	if err != nil {
		panic(err)
	}
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})
	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos:           repositories,
		Hasher:          hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:         cfg.JWT.SignKey,
		TokenTTL:        cfg.JWT.TokenTTL,
		StorageProvider: storage.NewFileStorage(uploader, cfg.S3.AwsBucket),
	}
	services := service.NewServices(deps)

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	// setup handler custom-validator as lib custom-validator
	handler.Validator = custom_validator.NewCustomValidator()
	v1.NewRouter(handler, services)

	// HTTP server
	log.Info("Starting http server...")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
