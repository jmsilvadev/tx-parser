package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestHappyPath(t *testing.T) {
	l := New(zapcore.DebugLevel)
	require.NotEmpty(t, l)
}
