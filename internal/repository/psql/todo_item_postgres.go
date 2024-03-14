package psql

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cheasezz/goTodo/internal/core"
	"github.com/Cheasezz/goTodo/pkg/postgres"
)

type TodoItem struct {
	db *postgres.Postgres
}

func NewTodoItemPostgres(db *postgres.Postgres) *TodoItem {
	return &TodoItem{db: db}
}

func (r *TodoItem) Create(ctx context.Context, listId int, item core.TodoItem) (int, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(ctx, createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)
	_, err = tx.Exec(ctx, createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	return itemId, tx.Commit(ctx)
}

func (r *TodoItem) GetAll(ctx context.Context, userId, listId int) ([]core.TodoItem, error) {
	var items []core.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
		INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Scany.Select(ctx, r.db.Pool, &items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItem) GetById(ctx context.Context, userId, itemId int) (core.TodoItem, error) {
	var item core.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
		INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Scany.Get(ctx, r.db.Pool, &item, query, itemId, userId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItem) Delete(ctx context.Context, userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul
		WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Pool.Exec(ctx, query, userId, itemId)
	return err
}

func (r *TodoItem) Update(ctx context.Context, userId, itemId int, input core.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
		WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)

	args = append(args, userId, itemId)

	_, err := r.db.Pool.Exec(ctx, query, args...)

	return err
}
