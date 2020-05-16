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
	AddTodo(context.Context, *todos.AddTodoRequest) (*empty.Empty, error)
	GetTodo(context.Context, uint64) (*todos.TodoResponse, error)
	GetTodos(context.Context, *empty.Empty) (*todos.TodosResponse, error)
	UpdateTodo(context.Context, *todos.UpdateTodoRequest) (*todos.TodoResponse, error)
	UpdateTodos(context.Context, *todos.UpdateTodosRequest) (*todos.TodosResponse, error)
	DeleteTodo(context.Context, uint64) (*empty.Empty, error)
	DeleteTodos(context.Context, []uint64) (*empty.Empty, error)
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
func (c *Conn) GetTodo(ctx context.Context, ID uint64) (*todos.TodoResponse, error) {
	log.Debugf("GetTodo() - ctx: %+v, id: %d", ctx, ID)

	t := &Todo{}

	query := `
SELECT t.todo_id, t.title, t.description, t.created_dt, t.updated_dt, ts.status
FROM app.todo t
JOIN app.todo_status ts
    ON t.status_id = ts.status_id
WHERE todo_id = %d
`

	if err := c.DB.Get(t, fmt.Sprintf(query, ID)); err != nil {
		log.WithError(err).Errorf("failed to retrieve Todo with id %d", ID)

		return nil, fmt.Errorf("failed to retrieve Todo with id %d", ID)
	}

	return newTodoResponse(t), nil
}

// GetTodos ...
func (c *Conn) GetTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error) {
	errMsg := errors.New("failed to retrieve all todos")
	txName := "GetTodos"

	t := Todo{}

	query := `
SELECT t.todo_id, t.title, t.description, t.created_dt, t.updated_dt, ts.status
FROM app.todo t
JOIN app.todo_status ts
    ON t.status_id = ts.status_id
`
	rows, err := c.DB.Queryx(query)

	if err != nil {
		log.WithError(err).Error(ErrQuery(txName).Error())
		return nil, errMsg
	}

	resp := []*todos.TodoResponse{}

	for rows.Next() {
		err := rows.StructScan(&t)
		if err != nil {
			log.WithError(err).Error(ErrScan(txName).Error())

			return nil, errMsg
		}

		resp = append(resp, newTodoResponse(&t))
	}

	return &todos.TodosResponse{
		Todos: resp,
	}, nil
}

// AddTodo ...
func (c *Conn) AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*empty.Empty, error) {
	errMsg := errors.New("failed to store Todo")
	txName := "AddTodo"

	query := `INSERT INTO app.todo (title, description) VALUES ($1, $2)`
	tx, err := c.DB.Begin()

	if err != nil {
		log.WithError(err).Error(ErrBeginTransaction(txName).Error())
		return nil, errMsg
	}

	res, err := tx.Exec(query, req.GetTitle(), req.GetDescription())

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName).Error())
		return nil, errMsg
	}

	if numRows, _ := res.RowsAffected(); numRows == 0 {
		log.WithError(err).Error(ErrNoRowsAffect(txName).Error())
		return nil, errMsg
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error(ErrCommit(txName).Error())
		return nil, errMsg
	}

	return &empty.Empty{}, nil
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
func (c *Conn) DeleteTodo(ctx context.Context, ID uint64) (*empty.Empty, error) {
	errMsg := errors.New("failed to delete Todo")
	txName := "DeleteTodo"

	query := `DELETE FROM app.todo WHERE todo_id = %d`
	tx, err := c.DB.Begin()

	if err != nil {
		log.WithError(err).Error(ErrBeginTransaction(txName).Error())
		return nil, errMsg
	}

	_, err = tx.Exec(fmt.Sprintf(query, ID))

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName).Error())
		return nil, errMsg
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error(ErrCommit(txName).Error())
		return nil, errMsg
	}

	return &empty.Empty{}, nil
}

// DeleteTodos ...
func (c *Conn) DeleteTodos(ctx context.Context, IDs []uint64) (*empty.Empty, error) {
	if len(IDs) == 0 {
		return &empty.Empty{}, nil
	}

	errMsg := errors.New("failed to delete Todos")
	txName := "DeleteTodos"

	arg := map[string]interface{}{"IDs": IDs}

	// dynamically bind the Todo IDs within the IN clause
	query, args, err := sqlx.Named("DELETE FROM app.todo WHERE todo_id IN (:IDs)", arg)

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName).Error())
		return nil, errMsg
	}

	query, args, err = sqlx.In(query, args...)

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName).Error())
		return nil, errMsg
	}

	query = c.DB.Rebind(query)

	tx, err := c.DB.Begin()

	if err != nil {
		log.WithError(err).Error(ErrBeginTransaction(txName).Error())
		return nil, errMsg
	}

	_, err = tx.Exec(query, args...)

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName).Error())
		return nil, errMsg
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error(ErrCommit(txName).Error())
		return nil, errMsg
	}

	return &empty.Empty{}, nil
}
