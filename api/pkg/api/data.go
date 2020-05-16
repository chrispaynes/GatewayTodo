package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chrispaynes/vorChall/proto/go/api/v1/todos"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// Data ...
type Data interface {
	AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*todos.TodoResponse, error)
	GetTodo(ctx context.Context, id uint64) (*todos.TodoResponse, error)
	GetTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error)
	UpdateTodo(ctx context.Context, req *todos.UpdateTodoRequest) (*todos.TodoResponse, error)
	UpdateTodos(ctx context.Context, req *todos.UpdateTodosRequest) (*todos.TodosResponse, error)
	DeleteTodo(ctx context.Context, req *todos.DeleteTodoRequest) (*empty.Empty, error)
	DeleteTodos(ctx context.Context, req *todos.DeleteTodosRequest) (*empty.Empty, error)
}

// Conn ...
type Conn struct {
	DB *sqlx.DB
}

// Todo ...
type Todo struct {
	ID          uint64         `db:"todo_id"`
	Title       sql.NullString `db:"title"`
	Description sql.NullString `db:"description"`
	Status      string         `db:"status"`
	CreatedDT   sql.NullString `db:"created_dt"`
	UpdatedDT   sql.NullString `db:"updated_dt"`
}

// NSString ...
func NSString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}

	return ns.String
}

func newTodoResponse(t *Todo) *todos.TodoResponse {
	return &todos.TodoResponse{
		Id:          t.ID,
		Title:       NSString(t.Title),
		Description: NSString(t.Description),
		CreatedDT:   NSString(t.CreatedDT),
		UpdatedDT:   NSString(t.UpdatedDT),
		Status:      t.Status,
	}
}

// GetTodo ...
func (c *Conn) GetTodo(ctx context.Context, id uint64) (*todos.TodoResponse, error) {
	log.Debugf("GetTodo() - ctx: %+v, id: %d", ctx, id)

	t := &Todo{}

	query := `
SELECT t.todo_id, t.title, t.description, t.created_dt, t.updated_dt, ts.status
FROM app.todo t
JOIN app.todo_status ts
    ON t.status_id = ts.status_id
WHERE todo_id = %d
`

	if err := c.DB.Get(t, fmt.Sprintf(query, id)); err != nil {
		log.WithError(err).Errorf("failed to retrieve Todo with id %d", id)

		return nil, fmt.Errorf("failed to retrieve Todo with id %d", id)
	}

	return newTodoResponse(t), nil
}

// GetTodos ...
func (c *Conn) GetTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error) {
	log.Debugf("GetTodos() - ctx: %+v", ctx)

	t := Todo{}

	query := `
SELECT t.todo_id, t.title, t.description, t.created_dt, t.updated_dt, ts.status
FROM app.todo t
JOIN app.todo_status ts
    ON t.status_id = ts.status_id
`
	rows, err := c.DB.Queryx(query)

	if err != nil {
		log.WithError(err).Error("failed to retrieve all Todos")
		return nil, errors.New("failed to retrieve all Todos")
	}

	resp := []*todos.TodoResponse{}

	for rows.Next() {
		err := rows.StructScan(&t)
		if err != nil {
			log.WithError(err).WithError(ErrScan("GetTodos"))

			return nil, errors.New("failed to retrieve all todos")
		}

		resp = append(resp, newTodoResponse(&t))
	}

	return &todos.TodosResponse{
		Todos: resp,
	}, nil
}

// AddTodo ...
func (c *Conn) AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*todos.TodoResponse, error) {
	panic("not implemented") // TODO: Implement
}

// UpdateTodo ...
func (c *Conn) UpdateTodo(ctx context.Context, req *todos.UpdateTodoRequest) (*todos.TodoResponse, error) {
	panic("not implemented") // TODO: Implement
}

// UpdateTodos ...
func (c *Conn) UpdateTodos(ctx context.Context, req *todos.UpdateTodosRequest) (*todos.TodosResponse, error) {
	panic("not implemented") // TODO: Implement
}

// DeleteTodo ...
func (c *Conn) DeleteTodo(ctx context.Context, req *todos.DeleteTodoRequest) (*empty.Empty, error) {
	panic("not implemented") // TODO: Implement
}

// DeleteTodos ...
func (c *Conn) DeleteTodos(ctx context.Context, req *todos.DeleteTodosRequest) (*empty.Empty, error) {
	panic("not implemented") // TODO: Implement
}
