package api

import (
	"github.com/chrispaynes/vorChall/pkg/logger"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type config struct {
	LogLevel string `default:"info" envconfig:"LOG_LEVEL" split_words:"true"`
}

var conf = &config{}

func init() {
	log.SetFormatter(&log.JSONFormatter{
		DisableTimestamp: true,
	})

	if err := envconfig.Process("", conf); err != nil {
		log.WithError(err).Fatal(logger.ErrParseEnv)
	}

	log.SetLevel(log.InfoLevel)

	if level, err := log.ParseLevel(conf.LogLevel); err != nil {
		log.SetLevel(level)
	}

	logger.LogEnv(conf, "api")
}
