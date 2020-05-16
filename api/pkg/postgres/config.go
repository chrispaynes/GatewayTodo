package postgres

import (
	"github.com/chrispaynes/vorChall/pkg/logger"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type config struct {
	DBHost     string `required:"true" envconfig:"POSTGRES_HOST" split_words:"true"`
	DBName     string `required:"true" envconfig:"POSTGRES_DB" split_words:"true"`
	DBPassword string `required:"true" envconfig:"POSTGRES_PASSWORD" split_words:"true"`
	DBUser     string `required:"true" envconfig:"POSTGRES_USER" split_words:"true"`
	LogLevel   string `default:"info" envconfig:"LOG_LEVEL" split_words:"true"`
}

var conf = &config{}

func init() {
	// TODO: remove these
	// os.Setenv("POSTGRES_DB", "")
	// os.Setenv("POSTGRES_PASSWORD", "")
	// os.Setenv("POSTGRES_USER", "")
	// os.Setenv("LOG_LEVEL", "info")

	// log.Printf
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

	logger.LogEnv(conf, "postgres")
}
