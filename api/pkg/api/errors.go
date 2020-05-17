package api

import (
	"errors"
	"fmt"
)

var (
	ErrBeginTransaction = func(queryName string) error {
		return newError(queryName, "failed to begin transaction")
	}

	ErrCommit = func(queryName string) error {
		return newError(queryName, "failed to commit transaction")
	}

	ErrDelete = func(queryName string) error {
		return newError(queryName, "failed to delete row(s)")
	}

	ErrExecTransaction = func(queryName string) error {
		return newError(queryName, "failed to execute transaction")
	}

	ErrInsert = func(queryName string) error {
		return newError(queryName, "failed to insert row(s)")
	}

	ErrNoRowsAffected = func(queryName string) error {
		return newError(queryName, "failed to affect any rows with transaction")
	}

	ErrQuery = func(queryName string) error {
		return newError(queryName, "failed to query row(s)")
	}

	InfoRollback = func(queryName string) error {
		return newError(queryName, "rolling back transaction")
	}

	ErrRollback = func(queryName string) error {
		return newError(queryName, "failed to rollback transaction")
	}

	ErrScan = func(queryName string) error {
		return newError(queryName, "failed to scan rows to destination")
	}
)

func newError(queryName, msg string) error {
	if queryName != "" {
		return fmt.Errorf("%s for %s", msg, queryName)
	}

	return errors.New(msg)
}
