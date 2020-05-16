package logger

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// ErrParseEnv ...
var ErrParseEnv = errors.New("failed to parse environment variables")

// SetFormat ...
func SetFormat(logLevel string) {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})

	log.SetLevel(log.InfoLevel)

	if level, err := log.ParseLevel(logLevel); err != nil {
		log.SetLevel(level)
	}

	return
}

// LogEnv ...
func LogEnv(conf interface{}, pkg string) {
	log.WithField("env", fmt.Sprintf("%+v", conf)).Infof("%s env", pkg)
}
