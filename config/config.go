package config

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type (
	Config struct {
		App    // `yaml:"app"`
		HTTP   //`yaml:"http"`
		Log    //`yaml:"log"`
		PG     //`yaml:"postgres"`
		JWT    //`yaml:"jwt"`
		Hasher //`yaml:"hasher"`
		S3
		//WebAPI `yaml:"webapi"`
	}

	App struct {
		Name    string //`env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string //`env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string //`env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string //`env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		MaxPoolSize int    // `env-required:"true" yaml:"max_pool_size" env:"PG_MAX_POOL_SIZE"`
		URL         string //`env-required:"false" env:"PG_URL"`
	}

	JWT struct {
		SignKey  string        //  `env-required:"true" env:"JWT_SIGN_KEY"`
		TokenTTL time.Duration //`env-required:"true" yaml:"token_ttl" env:"JWT_TOKEN_TTL"`
	}

	Hasher struct {
		Salt string //`env-required:"true" env:"HASHER_SALT"`
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
		os.Getenv("DB_HOST"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")) // os.Getenv("PG_URL") //"postgres://postgres:2001@localhost:5432/postgres?sslmode=disable"
	cfg.App.Name = "diplom-work"
	cfg.App.Version = "1.0.0"
	cfg.HTTP.Port = os.Getenv("HTTP_PORT") // "8080"
	cfg.Log.Level = "debug"
	cfg.JWT.SignKey = os.Getenv("JWT_SIGN_KEY") // "Diyorbek2001"
	cfg.JWT.TokenTTL = time.Minute * 120

	cfg.Hasher.Salt = os.Getenv("HASHER_SALT")                     // "Tatu65019$"
	cfg.S3.AwsRegion = os.Getenv("AWS_REGION")                     // "us-east-1"
	cfg.S3.AwsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")         //"AKIAZXWUWEZM54NWZVX2"
	cfg.S3.AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY") //"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"
	cfg.S3.AwsBucket = os.Getenv("AWS_BUCKET")                     //"test-go-diyorbek"
	os.Clearenv()

	return cfg, nil
}
