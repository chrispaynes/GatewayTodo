package logger

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// ErrParseEnv ...
var ErrParseEnv = errors.New("failed to parse environment variables")

// LogEnv ...
func LogEnv(conf interface{}, pkg string) {
	log.WithField("env", fmt.Sprintf("%+v", conf)).Infof("%s env", pkg)
}
