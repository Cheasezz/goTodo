package psql

import (
	"context"
	"fmt"

	"github.com/Cheasezz/goTodo/internal/core"
	"github.com/Cheasezz/goTodo/pkg/postgres"
	"github.com/sirupsen/logrus"
)

type Auth struct {
	db *postgres.Postgres
}

func NewAuthPostgres(db *postgres.Postgres) *Auth {
	return &Auth{db: db}
}

func (r *Auth) CreateUser(ctx context.Context, user core.User) (int, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}

	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", userTable)
	row := tx.QueryRow(ctx, query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	createUserSessionQuery := fmt.Sprintf("INSERT INTO %s (user_id) values ($1)", userSessionTable)
	_, err = tx.Exec(ctx, createUserSessionQuery, id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, nil
	}

	return id, tx.Commit(ctx)
}

func (r *Auth) GetUser(ctx context.Context, username, password string) (int, error) {
	var userId int
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &userId, query, username, password)
	
	return userId, err
}

func (r *Auth) SetSession(ctx context.Context, userId string, session core.Session) error {
	query := fmt.Sprintf("UPDATE %s us SET (refresh_token, expires_at) = ($1, $2) WHERE user_id = $3", userSessionTable)
	_, err := r.db.Pool.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, userId)

	return err
}

func (r *Auth) GetByRefreshToken(ctx context.Context, refreshToken string) (int, error) {
	var userId int

	query := fmt.Sprintf("SELECT user_id FROM %s WHERE refresh_token=$1", userSessionTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &userId, query, refreshToken)
	logrus.Printf("From GetByRefrashToken userId: %d", userId)
	return userId, err
}
