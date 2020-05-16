package server

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/chrispaynes/vorChall/pkg/api"
	"github.com/chrispaynes/vorChall/proto/go/api/v1/todos"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

/* TODOS
   - Add 12 factor for server ports
*/

// Config ...
type Config struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	grpcServer *grpc.Server
	httpServer *http.Server
	api.Todo
}

// NewServer ...
func NewServer(ctx context.Context, wg *sync.WaitGroup) *Config {
	return &Config{
		ctx: ctx,
		wg:  wg,
	}
}

// Start starts server
func (c *Config) Start() {
	c.wg.Add(1)
	go func() {
		log.Fatal(c.startGRPC())
		c.wg.Done()
	}()

	c.wg.Add(1)
	go func() {
		log.Fatal(c.startREST())
		c.wg.Done()
	}()

	c.wg.Wait()

	go func() {
		<-c.ctx.Done()
		log.Info("shutting down traffic to servers")
		go func() {
			log.Info("shutting down gRPC server")
			c.grpcServer.GracefulStop()
		}()

		go func() {
			log.Info("sending shutdown signal to the HTTP server")
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15*time.Second))

			defer cancel()

			if err := c.httpServer.Shutdown(ctx); err != nil {
			}
		}()
	}()
}

// Shutdown stops the server
func (c *Config) Shutdown() {
	c.cancel()
	c.wg.Wait()
}

func (c *Config) startGRPC() error {
	lis, err := net.Listen("tcp", ":3001")

	if err != nil {
		return err
	}

	c.grpcServer = grpc.NewServer()

	todos.RegisterTodosAPIServer(c.grpcServer, &api.TodoService{})

	c.grpcServer.Serve(lis)
	return nil
}

func (c *Config) startREST() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	c.httpServer = &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := todos.RegisterTodosAPIHandlerFromEndpoint(ctx, mux, ":3001", opts)
	if err != nil {
		return err
	}

	return c.httpServer.ListenAndServe()
}
