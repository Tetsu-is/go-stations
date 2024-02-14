package service

import (
	"context"
	"database/sql"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
// dbをメンバに持つTODOServiceを新しく作って返しているんだね
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
// 引数はcontext, subject, descriptionを受け取っているね
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	s.db.PrepareContext(ctx, insert)
	s.db.PrepareContext(ctx, confirm)

	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	row := s.db.QueryRowContext(ctx, confirm, id)

	var todo model.TODO
	row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	todo.ID = id

	if todo.Subject == "" {
		return &todo, nil
	}

	if todo.Description == "" {
		return &todo, nil
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	s.db.PrepareContext(ctx, update)
	s.db.PrepareContext(ctx, confirm)

	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	if count, _ := result.RowsAffected(); count == 0 {
		return nil, &model.ErrNotFound{}
	}

	row := s.db.QueryRowContext(ctx, confirm, id)

	var todo model.TODO
	row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	todo.ID = id

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
