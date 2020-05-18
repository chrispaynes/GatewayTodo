package server

import (
	"github.com/chrispaynes/vorChall/pkg/logger"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type config struct {
	gRPCPort string `default:"3001" envconfig:"GRPC_PORT" split_words:"true"`
	// this port is dynamically generated and set as PORT by Heroku
	HerokuPort         string `required:"true" envconfig:"PORT" split_words:"true"`
	HTTPTimeoutSeconds int    `default:"120" envconfig:"HTTP_TIME_SECONDS" split_words:"true"`
	LogLevel           string `default:"info" envconfig:"LOG_LEVEL" split_words:"true"`
	RESTPort           string `default:"3000" envconfig:"REST_POR" split_words:"true"`
	ServerEnv          string `default:"prod" envconfig:"SERVER_ENV" split_words:"true"`
	UIFilePath         string `default:"./dist/todo-app/" envconfig:"UI_FILEPATH" split_words:"true"`
	UIPort             string `default:"4200" envconfig:"UI_PORT" split_words:"true"`
}

var conf = &config{}

func init() {
	if err := envconfig.Process("", conf); err != nil {
		log.WithError(err).Fatal(logger.ErrParseEnv)
	}

	if conf.ServerEnv == "prod" {
		conf.RESTPort = conf.HerokuPort
	}

	if conf.ServerEnv == "prod" {
		conf.UIPort = conf.HerokuPort
	}

	logger.SetFormat(conf.LogLevel)
	logger.LogEnv(conf, "server")
}
