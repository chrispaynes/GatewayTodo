package api

import (
	"context"

	"github.com/chrispaynes/vorChall/proto/go/api/v1/todos"
	"github.com/golang/protobuf/ptypes/empty"
)

// TodoService represents a configuration
// link between the API and Data layers
type TodoService struct {
	Data Data
	API  Todo
}

// API interface represents a collection of funcs that satisfy gRPC route handlers generated from protobuf files
type API interface {
	AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*todos.TodoResponse, error)
	GetTodo(ctx context.Context, req *todos.GetTodoRequest) (*todos.TodoResponse, error)
	GetAllTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error)
	GetTodosById(ctx context.Context, req *todos.GetTodosRequest) (*todos.TodosResponse, error)
	UpdateTodo(ctx context.Context, req *todos.UpdateTodoRequest) (*todos.TodoResponse, error)
	UpdateTodos(ctx context.Context, req *todos.UpdateTodosRequest) (*todos.TodosResponse, error)
	DeleteTodo(ctx context.Context, req *todos.DeleteTodoRequest) (*empty.Empty, error)
	DeleteTodos(ctx context.Context, req *todos.DeleteTodosRequest) (*empty.Empty, error)
}

// AddTodo forwards a AddTodo request to the data layer
func (t *TodoService) AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*todos.TodoResponse, error) {
	return t.Data.AddTodo(ctx, req)
}

// GetTodo forwards a GetTodo request to the data layer
func (t *TodoService) GetTodo(ctx context.Context, req *todos.GetTodoRequest) (*todos.TodoResponse, error) {
	return t.Data.GetTodo(ctx, req.GetId())
}

// GetAllTodos forwards a GetAllTodos request to the data layer
func (t *TodoService) GetAllTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error) {
	return t.Data.GetAllTodos(ctx, nil)
}

// GetTodosById forwards a GetTodosById request to the data layer
func (t *TodoService) GetTodosById(ctx context.Context, req *todos.GetTodosRequest) (*todos.TodosResponse, error) {
	return t.Data.GetTodosByID(ctx, req)
}

// UpdateTodo forwards a UpdateTodo request to the data layer
func (t *TodoService) UpdateTodo(ctx context.Context, req *todos.UpdateTodoRequest) (*todos.TodoResponse, error) {
	return t.Data.UpdateTodo(ctx, req)
}

// UpdateTodos forwards a UpdateTodos request to the data layer
func (t *TodoService) UpdateTodos(ctx context.Context, req *todos.UpdateTodosRequest) (*todos.TodosResponse, error) {
	return t.Data.UpdateTodos(ctx, req)
}

// DeleteTodo forwards a DeleteTodo request to the data layer
func (t *TodoService) DeleteTodo(ctx context.Context, req *todos.DeleteTodoRequest) (*empty.Empty, error) {
	return t.Data.DeleteTodo(ctx, req.GetId())
}

// DeleteTodos forwards a DeleteTodos request to the data layer
func (t *TodoService) DeleteTodos(ctx context.Context, req *todos.DeleteTodosRequest) (*empty.Empty, error) {
	return t.Data.DeleteTodos(ctx, req.GetIds())
}
