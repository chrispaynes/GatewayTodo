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

var conf = &config{}

func init() {
	if err := envconfig.Process("", conf); err != nil {
		log.WithError(err).Fatal(logger.ErrParseEnv)
	}

	logger.SetFormat(conf.LogLevel)
	logger.LogEnv(conf, "main")
}

func main() {
	conf.ctx, conf.cancel = context.WithCancel(context.Background())

	db := postgres.NewDBWithRetry(1 * time.Minute)

	s := server.NewServer(conf.ctx, &sync.WaitGroup{}, db)

	s.Start()

	<-conf.ctx.Done()
	s.Shutdown()
}
