package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type (
	Config struct {
		App
		HTTP
		Log
		PG
		JWT
		Hasher
		S3
	}

	App struct {
		Name    string
		Version string
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}
	PG struct {
		MaxPoolSize int
		URL         string
	}

	JWT struct {
		SignKey  string
		TokenTTL time.Duration
	}

	Hasher struct {
		Salt string
	}

	S3 struct {
		AwsRegion, AwsAccessKeyId, AwsSecretAccessKey, AwsBucket string
	}
)

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
	log.Println(godotenv.Read())
	cfg := &Config{}
	cfg.MaxPoolSize = 20
	cfg.PG.URL = fmt.Sprintf("user=%s host=%s password=%s port=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	cfg.App.Name = "diplom-work"
	cfg.App.Version = "1.0.0"
	cfg.HTTP.Port = os.Getenv("HTTP_PORT")
	cfg.Log.Level = "debug"
	cfg.JWT.SignKey = os.Getenv("JWT_SIGN_KEY")
	cfg.JWT.TokenTTL = time.Minute * 120

	cfg.Hasher.Salt = os.Getenv("HASHER_SALT")
	cfg.S3.AwsRegion = os.Getenv("AWS_REGION")
	cfg.S3.AwsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	cfg.S3.AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	cfg.S3.AwsBucket = os.Getenv("AWS_BUCKET")
	os.Clearenv()

	return cfg, nil
}
