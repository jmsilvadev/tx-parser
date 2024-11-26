package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jmsilvadev/tx-parser/internal/handlers"
	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

type MockParser struct{}

func (m *MockParser) GetCurrentBlock() int {
	return 123
}

func (m *MockParser) Subscribe(address string) bool {
	return true
}

func (m *MockParser) GetTransactions(address string) []parser.Transaction {
	return []parser.Transaction{
		{
			Hash:        "0xabc",
			From:        "0x123",
			To:          "0x456",
			Value:       "100",
			BlockNumber: 1,
		},
	}
}

func (m *MockParser) UpdateBlockNumber() {}

func TestServer(t *testing.T) {
	mockParser := &MockParser{}
	mockLogger := logger.New(zapcore.DebugLevel)

	s := NewServer(
		func(s *Server) {
			s.port = ":8080"
			s.timeout = 10 * time.Second
			s.logger = mockLogger
			s.parser = mockParser
		},
	)

	go s.Start()
	time.Sleep(1 * time.Second)

	t.Run("HealthHandler", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.New(mockParser).HealthHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("GetCurrentBlock", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/get-current-block", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.New(mockParser).GetCurrentBlock)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "123")
	})

	t.Run("Subscribe", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/v1/subscribe?address=0x123", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.New(mockParser).Subscribe)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "true")
	})

	t.Run("GetTransactions", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/get-transactions?address=0x123", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.New(mockParser).GetTransactions)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "0xabc")
	})

	req, err := http.NewRequest("GET", "/shutdown", nil)
	assert.NoError(t, err)
	http.DefaultClient.Do(req)
}
