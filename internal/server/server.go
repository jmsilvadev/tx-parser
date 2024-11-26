package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jmsilvadev/tx-parser/internal/handlers"
	config "github.com/jmsilvadev/tx-parser/pkg/config"
	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
)

type Server struct {
	environment string
	port        string
	timeout     time.Duration
	logger      logger.Logger
	conf        *config.Config
	parser      parser.Parser
}

type ServerOption func(*Server)

func NewServer(options ...ServerOption) *Server {
	svr := &Server{}
	for _, opt := range options {
		opt(svr)
	}
	return svr
}

func (s *Server) Start() {
	go s.parser.UpdateBlockNumber()
	h := handlers.New(s.parser)

	http.HandleFunc("/health", h.HealthHandler)
	http.HandleFunc("/v1/get-current-block", h.GetCurrentBlock)
	http.HandleFunc("/v1/subscribe", h.Subscribe)
	http.HandleFunc("/v1/get-transactions", h.GetTransactions)
	http.HandleFunc("/", handlers.NotFoundHandler)

	server := &http.Server{
		Addr: s.port,
	}

	listener := make(chan os.Signal, 1)
	signal.Notify(listener, os.Interrupt, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s.logger.Warn(fmt.Sprint("received a shutdown signal:", <-listener))
		s.logger.Warn("shutdown the server...")
		server.Shutdown(context.Background())
		wg.Done()
	}()

	s.logger.Info("server listening at " + s.port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("failed to serve: " + err.Error())
	}

	wg.Wait()
	s.logger.Warn("server gracefully stopped")
}
