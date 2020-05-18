package logger

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var ErrParseEnv = errors.New("failed to parse environment variables")

// SetFormat sets the logger format within a package
func SetFormat(logLevel string) {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	})

	// output which files and line numbers printed the log
	// note: this is useful but can be pretty noisy
	log.SetReportCaller(true)

	log.SetLevel(log.InfoLevel)

	if level, err := log.ParseLevel(logLevel); err != nil {
		log.SetLevel(level)
	}

	return
}

// LogEnv logs the environment configuration to standard out
func LogEnv(conf interface{}, pkg string) {
	log.WithField("env", fmt.Sprintf("%+v", conf)).Infof("%s env", pkg)
}
