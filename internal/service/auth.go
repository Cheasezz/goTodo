package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Cheasezz/goTodo/internal/core"
	"github.com/Cheasezz/goTodo/pkg/auth"
	"github.com/Cheasezz/goTodo/pkg/hash"
)

const (
	tokenTTL = 12 * time.Hour
)

type AuthRepo interface {
	CreateUser(ctx context.Context, user core.User) (int, error)
	GetUser(ctx context.Context, username, password string) (core.User, error)
}

type Auth struct {
	repo         AuthRepo
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func newAuthService(r AuthRepo, h hash.PasswordHasher, tm auth.TokenManager) *Auth {
	return &Auth{
		repo:         r,
		hasher:       h,
		tokenManager: tm,
	}
}

func (s *Auth) CreateUser(ctx context.Context, user core.User) (int, error) {
	pass, err := s.hasher.Hash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = pass
	return s.repo.CreateUser(ctx, user)
}

func (s *Auth) GenerateToken(ctx context.Context, username, password string) (string, error) {
	pass, err := s.hasher.Hash(password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetUser(ctx, username, pass)
	if err != nil {
		return "", err
	}

	token, err := s.tokenManager.NewJWT(strconv.Itoa(user.Id), tokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Auth) ParseToken(accessToken string) (int, error) {
	claims, err := s.tokenManager.Parse(accessToken)
	if err != nil {
		return 0, err
	}

	intClaims, err := strconv.Atoi(claims)
	if err != nil {
		return 0, err
	}

	return intClaims, nil
}
