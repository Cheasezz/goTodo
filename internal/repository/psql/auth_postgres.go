package psql

import (
	"context"
	"fmt"

	"github.com/Cheasezz/goTodo/internal/core"
	"github.com/Cheasezz/goTodo/pkg/postgres"
)

type Auth struct {
	db *postgres.Postgres
}

func NewAuthPostgres(db *postgres.Postgres) *Auth {
	return &Auth{db: db}
}

func (r *Auth) CreateUser(ctx context.Context, user core.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", userTable)
	row := r.db.Pool.QueryRow(ctx, query, user.Name, user.Username, user.Password)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Auth) GetUser(ctx context.Context, username, password string) (core.User, error) {
	var user core.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &user, query, username, password)

	return user, err
}
