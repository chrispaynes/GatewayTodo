package api

import (
	"context"

	"github.com/chrispaynes/vorChall/proto/go/api/v1/todos"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
)

// TodoService ...
type TodoService struct {
	DB  *sqlx.DB
	API Todo
}

// Todo ...
type Todo interface {
	AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*todos.TodoResponse, error)
	GetTodo(ctx context.Context, req *todos.GetTodoRequest) (*todos.TodoResponse, error)
	GetTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error)
	UpdateTodo(ctx context.Context, req *todos.UpdateTodoRequest) (*todos.TodoResponse, error)
	UpdateTodos(ctx context.Context, req *todos.UpdateTodosRequest) (*todos.TodosResponse, error)
	DeleteTodo(ctx context.Context, req *todos.DeleteTodoRequest) (*empty.Empty, error)
	DeleteTodos(ctx context.Context, req *todos.DeleteTodosRequest) (*empty.Empty, error)
}

// AddTodo ...
func (t *TodoService) AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*todos.TodoResponse, error) {
	return &todos.TodoResponse{}, nil
}

// GetTodo ...
func (t *TodoService) GetTodo(ctx context.Context, req *todos.GetTodoRequest) (*todos.TodoResponse, error) {
	return &todos.TodoResponse{}, nil
}

// GetTodos ...
func (t *TodoService) GetTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error) {
	return &todos.TodosResponse{}, nil
}

// UpdateTodo ...
func (t *TodoService) UpdateTodo(ctx context.Context, req *todos.UpdateTodoRequest) (*todos.TodoResponse, error) {
	return &todos.TodoResponse{}, nil
}

// UpdateTodos ...
func (t *TodoService) UpdateTodos(ctx context.Context, req *todos.UpdateTodosRequest) (*todos.TodosResponse, error) {
	return &todos.TodosResponse{}, nil
}

// DeleteTodo ...
func (t *TodoService) DeleteTodo(ctx context.Context, req *todos.DeleteTodoRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// DeleteTodos ...
func (t *TodoService) DeleteTodos(ctx context.Context, req *todos.DeleteTodosRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
