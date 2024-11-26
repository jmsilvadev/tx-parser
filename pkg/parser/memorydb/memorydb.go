package memorydb

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmsilvadev/tx-parser/pkg/jsonrpc"
	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
)

var _ parser.Parser = &DB{}

type DB struct {
	currentBlock  int
	subscriptions map[string]bool
	transactions  map[string][]parser.Transaction
	jsonrpc       jsonrpc.JsonRpcClient
	logger        logger.Logger
	mu            sync.Mutex
}

func New(cli jsonrpc.JsonRpcClient, l logger.Logger) *DB {
	return &DB{
		subscriptions: make(map[string]bool),
		transactions:  make(map[string][]parser.Transaction),
		jsonrpc:       cli,
		logger:        l,
	}
}

func (p *DB) GetCurrentBlock(ctx context.Context) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentBlock
}

func (p *DB) Subscribe(ctx context.Context, address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.logger.Debug(strings.ToLower(address))

	if _, exists := p.subscriptions[strings.ToLower(address)]; exists {
		return false
	}

	p.subscriptions[strings.ToLower(address)] = true
	return true
}

func (p *DB) GetTransactions(ctx context.Context, address string) []parser.Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.transactions[strings.ToLower(address)]
}

func (p *DB) UpdateBlockNumber(ctx context.Context) {
	for {
		blockNumber, err := p.jsonrpc.GetCurrentBlockNumber(ctx)
		if err != nil {
			p.logger.Error(err.Error())
			continue
		}

		p.mu.Lock()
		if blockNumber > p.currentBlock {
			p.currentBlock = blockNumber
			transactions, err := p.jsonrpc.GetBlockTransactions(ctx, blockNumber)
			if err == nil {
				for _, tx := range transactions {
					if p.subscriptions[strings.ToLower(tx.From)] || p.subscriptions[strings.ToLower(tx.To)] {
						p.logger.Debug(fmt.Sprintf("%s | %s", tx.From, tx.To))
						p.transactions[strings.ToLower(tx.From)] = append(p.transactions[strings.ToLower(tx.From)], tx)
						p.transactions[strings.ToLower(tx.To)] = append(p.transactions[strings.ToLower(tx.To)], tx)
					}
				}
			} else {
				p.logger.Error(err.Error())
			}
		}
		p.mu.Unlock()

		time.Sleep(12 * time.Second)
	}
}
