package api

import (
	"errors"
	"fmt"
)

var (
	ErrScan = func(queryName string) error {
		msg := "failed to scan rows to destination"

		if queryName != "" {
			return fmt.Errorf("%s for %s", msg, queryName)
		}

		return errors.New(msg)
	}
)
