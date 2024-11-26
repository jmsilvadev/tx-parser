package leveldb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jmsilvadev/tx-parser/pkg/jsonrpc"
	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

var cliUrl = "https://ethereum-rpc.publicnode.com"

func setupTestDB(t *testing.T) *DB {
	l := logger.New(zapcore.DebugLevel)
	cli := jsonrpc.NewEthereum(l, cliUrl)
	path := "testdb"
	db, err := New(path, cli, l)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func teardownTestDB(db *DB) {
	db.db.Close()
	os.RemoveAll("testdb")
}

func TestUpdateBlockNumber(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// To cover lines
	go db.UpdateBlockNumber(context.Background())
	time.Sleep(time.Second)
}

func TestGetSetCurrentBlock(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Test setting and getting the current block
	err := db.SetCurrentBlock(context.Background(), 123)
	assert.NoError(t, err)

	block := db.GetCurrentBlock(context.Background())
	assert.Equal(t, 123, block)
}

func TestSubscribe(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Test subscribing an address
	subscribed := db.Subscribe(context.Background(), "0x123")
	assert.True(t, subscribed)

	// Test subscribing the same address again
	subscribed = db.Subscribe(context.Background(), "0x123")
	assert.False(t, subscribed)
}

func TestGetAddTransactions(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Test getting transactions for an address with no transactions
	txs := db.GetTransactions(context.Background(), "0x123")
	assert.Empty(t, txs)

	// Test adding a transaction
	tx := parser.Transaction{
		Hash:        "0xabc",
		From:        "0x123",
		To:          "0x456",
		Value:       "100",
		BlockNumber: 1,
	}
	err := db.AddTransaction(context.Background(), "0x123", tx)
	assert.NoError(t, err)

	// Test getting transactions for the address
	txs = db.GetTransactions(context.Background(), "0x123")
	assert.Len(t, txs, 1)
	assert.Equal(t, tx, txs[0])
}
