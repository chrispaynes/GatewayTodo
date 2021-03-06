package api

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chrispaynes/vorChall/proto/go/api/v1/todos"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// The Data interface serves provides functionality for the gRPC route handlers to
// access the the data layer and retrieve data needed for API responses
type Data interface {
	AddTodo(context.Context, *todos.AddTodoRequest) (*todos.TodoResponse, error)
	GetTodo(context.Context, uint32) (*todos.TodoResponse, error)
	GetAllTodos(context.Context, *empty.Empty) (*todos.TodosResponse, error)
	GetTodosByID(context.Context, *todos.GetTodosRequest) (*todos.TodosResponse, error)
	UpdateTodo(context.Context, *todos.UpdateTodoRequest) (*todos.TodoResponse, error)
	UpdateTodos(context.Context, *todos.UpdateTodosRequest) (*todos.TodosResponse, error)
	DeleteTodo(context.Context, uint32) (*empty.Empty, error)
	DeleteTodos(context.Context, []uint32) (*empty.Empty, error)
}

// Conn represents a connection to persistent storage
type Conn struct {
	DB *sqlx.DB
}

// Todo represents a Todo item stored and retrieved from the database
type Todo struct {
	ID          uint32         `db:"todo_id"`
	Title       sql.NullString `db:"title"`
	Description sql.NullString `db:"description"`
	Status      string         `db:"status"`
	CreatedDT   sql.NullString `db:"created_dt"`
	UpdatedDT   sql.NullString `db:"updated_dt"`
}

// NSString is a helper function to return a sqlNullstrings value or an empty string
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

// GetTodo gets a single todo from the database using the todo's ID
func (c *Conn) GetTodo(ctx context.Context, ID uint32) (*todos.TodoResponse, error) {
	log.Debugf("GetTodo() - ctx: %+v, id: %d", ctx, ID)
	errMsg := fmt.Errorf("failed to get Todo: %d", ID)
	txName := "GetTodo"

	t := &Todo{}

	query := `
SELECT t.todo_id, t.title, t.description, t.created_dt, t.updated_dt, ts.status
FROM app.todo t
JOIN app.todo_status ts
    ON t.status_id = ts.status_id
WHERE todo_id = %d
`

	if err := c.DB.Get(t, fmt.Sprintf(query, ID)); err != nil {
		log.WithError(err).Error(ErrQuery(txName))

		return nil, errMsg
	}

	return newTodoResponse(t), nil
}

// getTodoByTitleAndDescription since this postgres driver doesn't return primary key values after executing a command,
// this hacky func is an attempt to fetch a newly insert todo using its title and description.
func (c *Conn) getTodoByTitleAndDescription(ctx context.Context, title, description string) (*todos.TodoResponse, error) {
	errMsg := fmt.Errorf("failed to get Todo by title: %s", title)
	txName := "GetTodoByTitleAndDescription"

	t := &Todo{}

	query := `
SELECT t.todo_id, t.title, t.description, t.created_dt, t.updated_dt, ts.status
FROM app.todo t
JOIN app.todo_status ts
    ON t.status_id = ts.status_id
WHERE t.title = '%s'
    AND t.description = '%s'
ORDER BY t.created_dt DESC
LIMIT 1
`

	if err := c.DB.Get(t, fmt.Sprintf(query, title, description)); err != nil {
		log.WithError(err).Error(ErrQuery(txName))

		return nil, errMsg
	}

	return newTodoResponse(t), nil
}

// GetAllTodos gets all the todos from the database
// TODO: this should be augmented with query params to limit how many rows are returned
func (c *Conn) GetAllTodos(ctx context.Context, req *empty.Empty) (*todos.TodosResponse, error) {
	errMsg := errors.New("failed to retrieve all todos")
	txName := "GetAllTodos"

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

// GetTodosByID gets a collection of todos within a list of supplied ID's
func (c *Conn) GetTodosByID(ctx context.Context, req *todos.GetTodosRequest) (*todos.TodosResponse, error) {
	errMsg := fmt.Errorf("failed to retrieve todos: %v", req.GetIds())
	txName := "GetTodosByID"

	t := Todo{}

	// note: there's likely a limit of how many elements or characters Postgres allows in an IN clause
	// it's probably worth checking the length of req.GetIds() to set a sane limit
	q := `
    SELECT t.todo_id, t.title, t.description, t.created_dt, t.updated_dt, ts.status
    FROM app.todo t
    JOIN app.todo_status ts
        ON t.status_id = ts.status_id
    WHERE t.todo_id IN (:IDs)
    `

	arg := map[string]interface{}{"IDs": req.GetIds()}

	query, args, err := sqlx.Named(q, arg)

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName))
		return nil, errMsg
	}

	query, args, err = sqlx.In(query, args...)

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName))
		return nil, errMsg
	}

	query = c.DB.Rebind(query)

	rows, err := c.DB.Queryx(query, args...)

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

// AddTodo adds a new todo item in the database and attempts to fetch and return the recently add todo
func (c *Conn) AddTodo(ctx context.Context, req *todos.AddTodoRequest) (*todos.TodoResponse, error) {
	errMsg := errors.New("failed to store Todo")
	txName := "AddTodo"

	query := `INSERT INTO app.todo (title, description) VALUES ($1, $2)`
	tx, err := c.DB.Begin()

	if err != nil {
		log.WithError(err).Error(ErrBeginTransaction(txName))
		return nil, errMsg
	}

	res, err := tx.Exec(query, req.GetTitle(), req.GetDescription())

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName))
		return nil, errMsg
	}

	if numRows, _ := res.RowsAffected(); numRows == 0 {
		log.WithError(err).Error(ErrNoRowsAffected(txName))
		return nil, errMsg
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error(ErrCommit(txName))
		return nil, errMsg
	}

	resp, err := c.getTodoByTitleAndDescription(ctx, req.GetTitle(), req.GetDescription())

	if err != nil {
		// we're ignoring the error here since we successfully stored the new todo
		// the client will recieve a 200 OK response but won't see the new todo
		// returned in the response
		return nil, nil
	}

	return resp, nil
}

// UpdateTodo updates a todo in the database
func (c *Conn) UpdateTodo(ctx context.Context, req *todos.UpdateTodoRequest) (*todos.TodoResponse, error) {
	errMsg := fmt.Errorf("failed to update Todo: %d", req.GetId())
	txName := "UpdateTodo"

	query := `
UPDATE app.todo
SET title = '%s',
    description = '%s',
    updated_dt = NOW(),
    status_id = (SELECT status_id FROM app.todo_status WHERE status = '%s' LIMIT 1)
WHERE todo_id = %d
`
	tx, err := c.DB.Begin()

	if err != nil {
		log.WithError(err).Error(ErrBeginTransaction(txName))
		return nil, errMsg
	}

	res, err := tx.Exec(fmt.Sprintf(query, req.GetTitle(), req.GetDescription(), req.GetStatus(), req.GetId()))

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName))
		return nil, errMsg
	}

	if numRows, _ := res.RowsAffected(); numRows == 0 {
		log.WithError(err).Error(ErrNoRowsAffected(txName))
		return nil, errors.Wrap(errMsg, "no rows found")
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error(ErrCommit(txName))
		return nil, errMsg
	}

	return c.GetTodo(ctx, req.GetId())
}

// UpdateTodos performs a transactional update on a batch of todos but rolls the transaction back upon error
func (c *Conn) UpdateTodos(ctx context.Context, req *todos.UpdateTodosRequest) (*todos.TodosResponse, error) {
	errMsg := errors.New("failed to update Todos")
	txName := "UpdateTodos"

	query := `
UPDATE app.todo
SET title = '%s',
    description = '%s',
    updated_dt = NOW(),
    status_id = (SELECT status_id FROM app.todo_status WHERE status = '%s' LIMIT 1)
WHERE todo_id = %d
`
	tx, err := c.DB.Begin()

	if err != nil {
		log.WithError(err).Error(ErrBeginTransaction(txName))
		return nil, errMsg
	}

	IDs := []uint32{}

	for _, t := range req.GetTodos() {
		res, err := tx.Exec(fmt.Sprintf(query, t.GetTitle(), t.GetDescription(), t.GetStatus(), t.GetId()))

		if err != nil {
			log.WithError(err).Error(ErrExecTransaction(txName))
			break
		}

		if err != nil {
			log.WithError(err).Error(ErrExecTransaction(txName))
			break
		}

		if numRows, _ := res.RowsAffected(); numRows == 0 {
			log.WithError(err).Warn(ErrNoRowsAffected(txName))
			continue
		}

		IDs = append(IDs, t.GetId())
	}

	if err != nil {
		log.WithError(err).Error(InfoRollback(txName))

		if err = tx.Rollback(); err != nil {
			log.WithError(err).Error(ErrRollback(txName))
			return nil, errMsg
		}
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error(ErrCommit(txName))
		return nil, errMsg
	}

	return c.GetTodosByID(ctx, &todos.GetTodosRequest{Ids: IDs})
}

// DeleteTodo deletes a todo from the database using its ID
func (c *Conn) DeleteTodo(ctx context.Context, ID uint32) (*empty.Empty, error) {
	errMsg := fmt.Errorf("failed to delete Todo: %d", ID)
	txName := "DeleteTodo"

	query := `DELETE FROM app.todo WHERE todo_id = %d`
	tx, err := c.DB.Begin()

	if err != nil {
		log.WithError(err).Error(ErrBeginTransaction(txName))
		return nil, errMsg
	}

	_, err = tx.Exec(fmt.Sprintf(query, ID))

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName))
		return nil, errMsg
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error(ErrCommit(txName))
		return nil, errMsg
	}

	return &empty.Empty{}, nil
}

// DeleteTodos deletes a batch of todos based on their ID
func (c *Conn) DeleteTodos(ctx context.Context, IDs []uint32) (*empty.Empty, error) {
	if len(IDs) == 0 {
		return &empty.Empty{}, nil
	}

	errMsg := fmt.Errorf("failed to delete Todos: %v", IDs)
	txName := "DeleteTodos"

	arg := map[string]interface{}{"IDs": IDs}

	// dynamically bind the Todo IDs within the IN clause
	query, args, err := sqlx.Named("DELETE FROM app.todo WHERE todo_id IN (:IDs)", arg)

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName))
		return nil, errMsg
	}

	query, args, err = sqlx.In(query, args...)

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName))
		return nil, errMsg
	}

	query = c.DB.Rebind(query)

	tx, err := c.DB.Begin()

	if err != nil {
		log.WithError(err).Error(ErrBeginTransaction(txName))
		return nil, errMsg
	}

	_, err = tx.Exec(query, args...)

	if err != nil {
		log.WithError(err).Error(ErrExecTransaction(txName))
		return nil, errMsg
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error(ErrCommit(txName))
		return nil, errMsg
	}

	return &empty.Empty{}, nil
}
