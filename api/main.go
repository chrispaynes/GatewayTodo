package main

import (
	"context"
	"sync"
	"time"

	"github.com/chrispaynes/vorChall/pkg/logger"
	"github.com/chrispaynes/vorChall/pkg/postgres"
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
	LogLevel string `default:"info" envconfig:"LOG_LEVEL" split_words:"true"`
	ctx      context.Context
	cancel   context.CancelFunc
}

var cfg = &config{}

func init() {
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

	db := postgres.NewDBWithRetry(1 * time.Minute)

	s := server.NewServer(cfg.ctx, &sync.WaitGroup{}, db)

	s.Start()

	<-cfg.ctx.Done()
	s.Shutdown()
}
