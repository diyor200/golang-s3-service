package service

import (
	"context"
	"golang-diplom-work/internal/entity"
	"golang-diplom-work/internal/repo"
	"golang-diplom-work/pkg/hasher"
	"golang-diplom-work/pkg/storage"
	"time"
)

type AuthCreateUserInput struct {
	Username, Password string
}

type AuthGenerateTokenInput struct {
	Username, Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (int, error)
}

// Files interface for s3
type Files interface {
	Upload(ctx context.Context, file entity.File) (string, error)
}

type Services struct {
	Auth
	Files
}

type ServicesDependencies struct {
	Repos           *repo.Repositories
	Hasher          hasher.PasswordHasher
	StorageProvider storage.Provider
	SignKey         string
	TokenTTL        time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:  NewAuthService(deps.Repos, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Files: NewFileService(deps.StorageProvider),
	}
}
