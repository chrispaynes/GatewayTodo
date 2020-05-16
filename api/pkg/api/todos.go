package api

import (
	"context"

	"github.com/chrispaynes/vorChall/proto/go/api/v1/todos"
	"github.com/golang/protobuf/ptypes/empty"
)

// TodoService ...
type TodoService struct {
	Data Data
	API  Todo
}

// API ...
type API interface {
	AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*empty.Empty, error)
	GetTodo(ctx context.Context, req *todos.GetTodoRequest) (*todos.TodoResponse, error)
	GetTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error)
	UpdateTodo(ctx context.Context, req *todos.UpdateTodoRequest) (*todos.TodoResponse, error)
	UpdateTodos(ctx context.Context, req *todos.UpdateTodosRequest) (*todos.TodosResponse, error)
	DeleteTodo(ctx context.Context, req *todos.DeleteTodoRequest) (*empty.Empty, error)
	DeleteTodos(ctx context.Context, req *todos.DeleteTodosRequest) (*empty.Empty, error)
}

// AddTodo ...
func (t *TodoService) AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*empty.Empty, error) {
	return t.Data.AddTodo(ctx, req)
}

// GetTodo ...
func (t *TodoService) GetTodo(ctx context.Context, req *todos.GetTodoRequest) (*todos.TodoResponse, error) {
	return t.Data.GetTodo(ctx, req.GetId())
}

// GetTodos ...
func (t *TodoService) GetTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error) {
	return t.Data.GetTodos(ctx, nil)
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
	return t.Data.DeleteTodo(ctx, req.GetId())
}

// DeleteTodos ...
func (t *TodoService) DeleteTodos(ctx context.Context, req *todos.DeleteTodosRequest) (*empty.Empty, error) {
	return t.Data.DeleteTodos(ctx, req.GetIds())
}
