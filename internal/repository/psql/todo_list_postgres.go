package psql

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cheasezz/goTodo/internal/core"
	"github.com/Cheasezz/goTodo/pkg/postgres"
)

type TodoList struct {
	db *postgres.Postgres
}

func NewTodoListPostgres(db *postgres.Postgres) *TodoList {
	return &TodoList{db: db}
}

func (r *TodoList) Create(ctx context.Context, userId int, list core.TodoList) (int, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(ctx, createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	createUserListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(ctx, createUserListQuery, userId, id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, nil
	}

	return id, tx.Commit(ctx)
}

func (r *TodoList) GetAll(ctx context.Context, userId int) ([]core.TodoList, error) {
	var lists []core.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := r.db.Scany.Select(ctx, r.db.Pool, &lists, query, userId)

	return lists, err
}

func (r *TodoList) GetById(ctx context.Context, userId, listId int) (core.TodoList, error) {
	var list core.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
			INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	err := r.db.Scany.Get(ctx, r.db.Pool, &list, query, userId, listId)

	return list, err
}

func (r *TodoList) Delete(ctx context.Context, userId, listId int) error {
	query := fmt.Sprintf(`DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2`,
		todoListsTable, usersListsTable)

	_, err := r.db.Pool.Exec(ctx, query, userId, listId)

	return err
}

func (r *TodoList) Update(ctx context.Context, userId, listId int, input core.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, listId, userId)


	_, err := r.db.Pool.Exec(ctx, query, args...)

	return err
}
