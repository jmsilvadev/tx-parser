package config

import (
	"context"
	"testing"
	"time"

	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewConfig(t *testing.T) {
	l := logger.New(zap.DebugLevel)
	parser, _ := getDatabase("memorydb", "", cliUrl, l)
	got := New(context.Background(), ":5000", "dev", time.Second, parser, &zap.Logger{})
	if got.ServerPort != ":5000" {
		t.Errorf("Got and Expected are not equals. Got: %v, expected: :5000", got.ServerPort)
	}
}

func TestGetDeaultConfig(t *testing.T) {
	config := GetDefaultConfig()
	if config.ServerPort == "" {
		t.Errorf("Got and Expected are not equals. got: '', expected: !''")
	}
}

func TestGetEnv(t *testing.T) {
	v := getEnv("a", "b")
	require.Equal(t, "b", v)
}
