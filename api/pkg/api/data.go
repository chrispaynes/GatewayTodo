package api

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrispaynes/vorChall/proto/go/api/v1/todos"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// Data ...
type Data interface {
	AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*todos.TodoResponse, error)
	GetTodo(ctx context.Context, id int64) (*todos.TodoResponse, error)
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
	ID          int64          `db:"id"`
	Title       sql.NullString `db:"title"`
	Description sql.NullString `db:"description"`
	Status      string         `db:"status"`
}

// NSString ...
func NSString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}

	return ns.String
}

// GetTodo ...
func (c *Conn) GetTodo(ctx context.Context, id int64) (*todos.TodoResponse, error) {
	log.Debugf("GetTodo() - ctx: %+v, id: %d", ctx, id)

	t := &Todo{}

	if err := c.DB.Get(t, fmt.Sprintf(`SELECT id, title, description, status FROM app.todo WHERE id = %d`, id)); err != nil {
		log.WithError(err).Errorf("failed to retrieve Todo with id %d", id)

		return nil, fmt.Errorf("failed to retrieve Todo with id %d", id)
	}

	return &todos.TodoResponse{
		Id:          t.ID,
		Title:       NSString(t.Title),
		Description: NSString(t.Description),
		Status:      t.Status,
	}, nil
}

// GetTodos ...
func (c *Conn) GetTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error) {
	log.Debugf("GetTodos() - ctx: %+v", ctx)

	t := Todo{}
	rows, err := c.DB.Queryx("SELECT id, title, description, status FROM app.todo ORDER BY id DESC")

	if err != nil {
		log.Fatalln(err)
	}

	resp := []*todos.TodoResponse{}

	for rows.Next() {
		err := rows.StructScan(&t)
		if err != nil {
			log.Fatalln(err)
		}

		resp = append(resp, &todos.TodoResponse{
			Id:          t.ID,
			Title:       NSString(t.Title),
			Description: NSString(t.Description),
			Status:      t.Status,
		})
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
