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
	tokenTTL = 10 * time.Minute
	refreshTokenTTL = 24 * time.Hour
)

type AuthRepo interface {
	CreateUser(ctx context.Context, user core.User) (int, error)
	GetUser(ctx context.Context, username, password string) (int, error)
	SetSession(ctx context.Context, userId string, session core.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (int, error)
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

func (s *Auth) SignIn(ctx context.Context, username, password string) (auth.Tokens, error) {
	pass, err := s.hasher.Hash(password)
	if err != nil {
		return auth.Tokens{}, err
	}

	userId, err := s.repo.GetUser(ctx, username, pass)
	if err != nil {
		return auth.Tokens{}, err
	}

	return s.createSession(ctx, strconv.Itoa(userId))
}

func (s *Auth) createSession(ctx context.Context, userId string) (auth.Tokens, error) {
	var (
		res auth.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, tokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := core.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userId, session)

	return res, err
}

func (s *Auth) RefreshTokens(ctx context.Context, refreshToken string) (auth.Tokens, error) {
	userId, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return auth.Tokens{}, err
	}

	return s.createSession(ctx, strconv.Itoa(userId))
}