package main

import (
	server "github.com/jmsilvadev/tx-parser/internal/server"
	"github.com/jmsilvadev/tx-parser/pkg/config"
)

func main() {
	c := config.GetDefaultConfig()
	run(c)
}

func run(conf *config.Config) error {
	serverOptions := []server.ServerOption{
		server.WithPort(conf.ServerPort),
		server.WithEnvironment(conf.Env),
		server.WithLogger(conf.Logger),
		server.WithParser(conf.Parser),
	}

	s := server.NewServer(serverOptions...)

	s.Start()

	return nil
}
