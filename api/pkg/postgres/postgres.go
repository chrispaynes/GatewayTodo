package postgres

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

// NewDBWithRetry ...
func NewDBWithRetry(t time.Duration) *sqlx.DB {
	log.Info("attempting to connect to Postgres")

	timeout := time.Now().Add(t)

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBName)

	attempt := 0

	for {
		attempt++

		db, err := sqlx.Connect("postgres", dsn)

		if err == nil {
			log.Info("successfully connected to Postgres")
			return db
		}

		if time.Now().After(timeout) {
			log.Fatal(errors.Wrapf(err, "failed to connect to postgres database after attempt %d", attempt))
		}

		time.Sleep(10 * time.Second)
	}

}
