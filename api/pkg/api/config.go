package api

import (
	"github.com/chrispaynes/vorChall/pkg/logger"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type config struct {
	LogLevel  string `default:"info" envconfig:"LOG_LEVEL" split_words:"true"`
	ServerEnv string `default:"prod" envconfig:"SERVER_ENV" split_words:"true"`
}

var conf = &config{}

func init() {
	if err := envconfig.Process("", conf); err != nil {
		log.WithError(err).Fatal(logger.ErrParseEnv)
	}

	logger.SetFormat(conf.LogLevel)
	logger.LogEnv(conf, "api")
}
