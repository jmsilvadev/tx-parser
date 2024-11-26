package server

import (
	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
)

func WithPort(v string) ServerOption {
	return func(s *Server) {
		s.port = v
	}
}

func WithEnvironment(v string) ServerOption {
	return func(s *Server) {
		s.environment = v
	}
}

func WithParser(v parser.Parser) ServerOption {
	return func(s *Server) {
		s.parser = v
	}
}

func WithLogger(v logger.Logger) ServerOption {
	return func(s *Server) {
		s.logger = v
	}
}
