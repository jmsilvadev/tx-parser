package memorydb

import (
	"testing"
	"time"

	"github.com/jmsilvadev/tx-parser/pkg/jsonrpc"
	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

var cliUrl = "https://ethereum-rpc.publicnode.com"

func TestUpdateBlockNumber(t *testing.T) {
	l := logger.New(zapcore.DebugLevel)
	cli := jsonrpc.NewEthereum(l, cliUrl)
	db := New(cli, l)

	// To cover lines
	go db.UpdateBlockNumber()
	time.Sleep(time.Second)
}

func TestGetCurrentBlock(t *testing.T) {
	l := logger.New(zapcore.DebugLevel)
	cli := jsonrpc.NewEthereum(l, cliUrl)
	db := New(cli, l)

	block := db.GetCurrentBlock()
	assert.Equal(t, 0, block)

	db.mu.Lock()
	db.currentBlock = 123
	db.mu.Unlock()

	block = db.GetCurrentBlock()
	assert.Equal(t, 123, block)
}

func TestSubscribe(t *testing.T) {
	l := logger.New(zapcore.DebugLevel)
	cli := jsonrpc.NewEthereum(l, cliUrl)
	db := New(cli, l)

	subscribed := db.Subscribe("0x123")
	assert.True(t, subscribed)

	subscribed = db.Subscribe("0x123")
	assert.False(t, subscribed)
}

func TestGetTransactions(t *testing.T) {
	l := logger.New(zapcore.DebugLevel)
	cli := jsonrpc.NewEthereum(l, cliUrl)
	db := New(cli, l)

	txs := db.GetTransactions("0x123")
	assert.Empty(t, txs)

	tx := parser.Transaction{
		Hash:        "0xabc",
		From:        "0x123",
		To:          "0x456",
		Value:       "100",
		BlockNumber: 1,
	}

	db.mu.Lock()
	db.transactions["0x123"] = append(db.transactions["0x123"], tx)
	db.mu.Unlock()

	txs = db.GetTransactions("0x123")
	assert.Len(t, txs, 1)
	assert.Equal(t, tx, txs[0])
}
