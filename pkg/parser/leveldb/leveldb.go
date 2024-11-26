package leveldb

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmsilvadev/tx-parser/pkg/jsonrpc"
	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var _ parser.Parser = &DB{}

type DB struct {
	db      *leveldb.DB
	jsonrpc jsonrpc.JsonRpcClient
	logger  logger.Logger
}

func New(path string, cli jsonrpc.JsonRpcClient, l logger.Logger) (*DB, error) {
	db, err := leveldb.OpenFile(path, &opt.Options{
		ErrorIfMissing: false,
	})
	if err != nil {
		return nil, err
	}
	return &DB{
		db:      db,
		jsonrpc: cli,
		logger:  l,
	}, nil
}

func (p *DB) GetCurrentBlock() int {
	data, err := p.db.Get([]byte("currentBlock"), nil)
	if err != nil {
		p.logger.Debug(err.Error())
		if err == leveldb.ErrNotFound {
			return 0
		}
		return 0
	}
	var block int
	if err := json.Unmarshal(data, &block); err != nil {
		p.logger.Debug(err.Error())
		return 0
	}
	return block
}

func (p *DB) SetCurrentBlock(block int) error {
	data, err := json.Marshal(block)
	if err != nil {
		p.logger.Debug(err.Error())
		return err
	}
	err = p.db.Put([]byte("currentBlock"), data, nil)
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}

func (p *DB) Subscribe(address string) bool {
	_, err := p.db.Get([]byte("subscribed:"+strings.ToLower(address)), nil)
	if err == nil {
		return false
	}

	if err != leveldb.ErrNotFound {
		return false
	}

	err = p.db.Put([]byte("subscribed:"+strings.ToLower(address)), []byte("true"), nil)
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err == nil
}

func (p *DB) GetTransactions(address string) []parser.Transaction {
	data, err := p.db.Get([]byte("transactions:"+strings.ToLower(address)), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return []parser.Transaction{}
		}
		p.logger.Debug(err.Error())
		return nil
	}
	var transactions []parser.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		p.logger.Debug(err.Error())
		return nil
	}
	return transactions
}

func (p *DB) AddTransaction(address string, tx parser.Transaction) error {
	transactions := p.GetTransactions(strings.ToLower(address))
	if len(transactions) == 0 {
		transactions = []parser.Transaction{}
	}

	transactions = append(transactions, tx)
	data, err := json.Marshal(transactions)
	if err != nil {
		p.logger.Debug(err.Error())
		return err
	}

	err = p.db.Put([]byte("transactions:"+strings.ToLower(address)), data, nil)
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}

func (p *DB) UpdateBlockNumber() {
	for {
		blockNumber, err := p.jsonrpc.GetCurrentBlockNumber()
		if err != nil {
			p.logger.Debug(err.Error())
			continue
		}

		currentBlock := p.GetCurrentBlock()
		if blockNumber > currentBlock {
			err = p.SetCurrentBlock(blockNumber)
			if err != nil {
				p.logger.Debug(err.Error())
				continue
			}
			transactions, err := p.jsonrpc.GetBlockTransactions(blockNumber)
			if err == nil {
				for _, tx := range transactions {
					subscribedFrom, _ := p.db.Get([]byte("subscribed:"+strings.ToLower(tx.From)), nil)
					subscribedTo, _ := p.db.Get([]byte("subscribed:"+strings.ToLower(tx.To)), nil)
					if subscribedFrom != nil || subscribedTo != nil {
						p.logger.Debug(fmt.Sprintf("AddTransaction address %s %s ", tx.From, tx.To))
						p.AddTransaction(strings.ToLower(tx.From), tx)
						p.AddTransaction(strings.ToLower(tx.To), tx)
					}
				}
			} else {
				p.logger.Debug(err.Error())
			}
		}

		time.Sleep(12 * time.Second)
	}
}
