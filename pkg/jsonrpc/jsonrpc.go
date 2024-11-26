package jsonrpc

import (
	"context"

	"github.com/jmsilvadev/tx-parser/pkg/parser"
)

type JsonRpcClient interface {
	GetCurrentBlockNumber(context.Context) (int, error)
	GetBlockTransactions(context.Context, int) ([]parser.Transaction, error)
}
