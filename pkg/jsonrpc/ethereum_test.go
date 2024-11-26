package jsonrpc

import (
	"context"
	"testing"

	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

var curBlock = 0
var cliUrl = "https://ethereum-rpc.publicnode.com"

func TestGetCurrentBlockNumber(t *testing.T) {
	l := logger.New(zapcore.DebugLevel)

	e := NewEthereum(l, cliUrl)
	blockNumber, err := e.GetCurrentBlockNumber(context.Background())
	assert.NoError(t, err)
	assert.Greater(t, blockNumber, 0)

	curBlock = blockNumber
}

func TestGetBlockTransactions(t *testing.T) {
	l := logger.New(zapcore.DebugLevel)
	e := NewEthereum(l, cliUrl)
	transactions, err := e.GetBlockTransactions(context.Background(), curBlock)
	assert.NoError(t, err)
	assert.Greater(t, len(transactions), 0)
}
