package postgres

import (
	"github.com/chrispaynes/vorChall/pkg/logger"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type config struct {
	DBHost            string `required:"true" envconfig:"POSTGRES_HOST" split_words:"true"`
	DBName            string `required:"true" envconfig:"POSTGRES_DB" split_words:"true"`
	DBPassword        string `required:"true" envconfig:"POSTGRES_PASSWORD" split_words:"true"`
	DBUser            string `required:"true" envconfig:"POSTGRES_USER" split_words:"true"`
	LogLevel          string `default:"info" envconfig:"LOG_LEVEL" split_words:"true"`
	RetrySleepSeconds int64  `default:"15" envconfig:"RETRY_SLEEP_SECONDS" split_words:"true"`
	EnableSSL         bool   `default:"true" envconfig:"ENABLE_SSL" split_words:"true"`
	ServerEnv         string `default:"prod" envconfig:"SERVER_ENV" split_words:"true"`
}

var conf = &config{}

func init() {
	if err := envconfig.Process("", conf); err != nil {
		log.WithError(err).Fatal(logger.ErrParseEnv)
	}

	logger.SetFormat(conf.LogLevel)

	// don't log in production because DB secrets will be leaked
	if conf.ServerEnv == "dev" {
		logger.LogEnv(conf, "postgres")
	}
}
