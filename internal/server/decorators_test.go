package server

import (
	"testing"

	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestWithPort(t *testing.T) {
	s := &Server{}
	opt := WithPort(":8080")
	opt(s)
	assert.Equal(t, ":8080", s.port)
}

func TestWithEnvironment(t *testing.T) {
	s := &Server{}
	opt := WithEnvironment("development")
	opt(s)
	assert.Equal(t, "development", s.environment)
}

func TestWithParser(t *testing.T) {
	s := &Server{}
	mockParser := &MockParser{}
	opt := WithParser(mockParser)
	opt(s)
	assert.Equal(t, mockParser, s.parser)
}

func TestWithLogger(t *testing.T) {
	s := &Server{}
	mockLogger := logger.New(zapcore.DebugLevel)
	opt := WithLogger(mockLogger)
	opt(s)
	assert.Equal(t, mockLogger, s.logger)
}
