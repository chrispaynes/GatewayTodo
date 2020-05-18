package server

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/chrispaynes/vorChall/pkg/api"
	"github.com/chrispaynes/vorChall/proto/go/api/v1/todos"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Config represents the server connections and resources
type Config struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	grpcServer *grpc.Server
	httpServer *http.Server
	api.Todo
	db *sqlx.DB
}

// NewServer creates a new server configuration
func NewServer(ctx context.Context, wg *sync.WaitGroup, db *sqlx.DB) *Config {
	return &Config{
		ctx: ctx,
		wg:  wg,
		db:  db,
	}
}

// Start starts the gRPC and REST servers and listens
// for a context.Done() signal before shutting down the servers
func (c *Config) Start() {
	c.wg.Add(1)
	go func() {
		log.Info("starting the gRPC server")
		log.Fatal(c.startGRPCServer())
		c.wg.Done()
	}()

	c.wg.Add(1)
	go func() {
		log.Info("starting the REST server")
		log.Fatal(c.startRESTServer())
		c.wg.Done()
	}()

	log.Info("successfully started gRPC and REST servers")
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

// Shutdown cancels the context and shuts down the the server
func (c *Config) Shutdown() {
	c.cancel()
	c.wg.Wait()
}

// startGRPCServer starts the gRPC server
// and registers the gRPC API Server
func (c *Config) startGRPCServer() error {
	lis, err := net.Listen("tcp", ":"+conf.GRPCPort)

	if err != nil {
		return err
	}

	c.grpcServer = grpc.NewServer()

	todos.RegisterTodosAPIServer(c.grpcServer, &api.TodoService{
		Data: &api.Conn{
			DB: c.db,
		},
	})

	c.grpcServer.Serve(lis)
	return nil
}

// startRESTServer starts the REST HTTP server
// and sets the gRPC Gateway mux as its handler
func (c *Config) startRESTServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		// use the original protobuf name for fields.
		OrigName: true,
		// render fields with zero values.
		EmitDefaults: true,
	})

	c.httpServer = &http.Server{
		Addr:         ":" + conf.RESTPort,
		IdleTimeout:  time.Duration(conf.HTTPTimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(conf.HTTPTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(conf.HTTPTimeoutSeconds) * time.Second,
		Handler:      newRouter(mux),
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := todos.RegisterTodosAPIHandlerFromEndpoint(ctx, mux, ":"+conf.GRPCPort, opts)

	if err != nil {
		return err
	}

	return c.httpServer.ListenAndServe()
}

// newRouter creates a new HTTP router with middleware
// and serves to proxy REST traffic to the gRPC Gateway mux
func newRouter(mux *runtime.ServeMux) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(time.Duration(15 * time.Second)))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// forward v1 todo traffic to the gRPC mux
	pattern := "/v1/todo*"
	r.Method("DELETE", pattern, mux)
	r.Method("GET", pattern, mux)
	r.Method("PUT", pattern, mux)
	r.Method("POST", pattern, mux)

	return r
}
