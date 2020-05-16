package main

import (
	"context"
	"os"
	"sync"

	"github.com/chrispaynes/vorChall/pkg/logger"
	"github.com/chrispaynes/vorChall/pkg/server"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	_ "github.com/golang/glog"
	_ "github.com/lib/pq"
)

/* TODOS
   - Add 12 factor
*/

type config struct {
	DBName     string `required:"true" envconfig:"POSTGRES_DB" split_words:"true"`
	DBPassword string `required:"true" envconfig:"POSTGRES_PASSWORD" split_words:"true"`
	DBUser     string `required:"true" envconfig:"POSTGRES_USER" split_words:"true"`
	LogLevel   string `default:"info" envconfig:"LOG_LEVEL" split_words:"true"`
	ctx        context.Context
	cancel     context.CancelFunc
}

var cfg = &config{}

func init() {
	// TODO: remove these
	os.Setenv("POSTGRES_DB", "")
	os.Setenv("POSTGRES_PASSWORD", "")
	os.Setenv("POSTGRES_USER", "")
	os.Setenv("LOG_LEVEL", "info")

	// log.Printf
	log.SetFormatter(&log.JSONFormatter{
		DisableTimestamp: true,
	})

	if err := envconfig.Process("", cfg); err != nil {
		log.WithError(err).Fatal(logger.ErrParseEnv)
	}

	log.SetLevel(log.InfoLevel)

	if level, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.SetLevel(level)
	}

	logger.LogEnv(cfg, "main")
}

func main() {
	cfg.ctx, cfg.cancel = context.WithCancel(context.Background())

	s := server.NewServer(cfg.ctx, &sync.WaitGroup{})

	s.Start()

	<-cfg.ctx.Done()
	s.Shutdown()

}
