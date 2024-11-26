package jsonrpc

import "github.com/jmsilvadev/tx-parser/pkg/parser"

type JsonRpcClient interface {
	GetCurrentBlockNumber() (int, error)
	GetBlockTransactions(blockNumber int) ([]parser.Transaction, error)
}
