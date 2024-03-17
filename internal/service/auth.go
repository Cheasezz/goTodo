package service

import (
	"context"
	"errors"
	"time"

	"github.com/Cheasezz/goTodo/internal/core"
	"github.com/Cheasezz/goTodo/pkg/hash"
	"github.com/dgrijalva/jwt-go"
)

const (
	signingKey = "qqvgfg5jk3fwioi#ifsd"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthRepo interface {
	CreateUser(ctx context.Context, user core.User) (int, error)
	GetUser(ctx context.Context, username, password string) (core.User, error)
}

type Auth struct {
	repo   AuthRepo
	hasher hash.PasswordHasher
}

func newAuthService(repo AuthRepo, hasher hash.PasswordHasher) *Auth {
	return &Auth{repo: repo, hasher: hasher}
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *Auth) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not type *tokenClaims")
	}
	return claims.UserId, nil
}
