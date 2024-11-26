package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
)

type Ethereum struct {
	log    logger.Logger
	cliUrl string
}

var _ JsonRpcClient = &Ethereum{}

func NewEthereum(l logger.Logger, cliUrl string) *Ethereum {
	return &Ethereum{
		log:    l,
		cliUrl: cliUrl,
	}
}

func (e *Ethereum) GetCurrentBlockNumber(ctx context.Context) (int, error) {
	e.log.Debug("Executing GetCurrentBlockNumber")
	payload := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`

	resp, err := http.Post(e.cliUrl, "application/json", strings.NewReader(payload))
	if err != nil {
		e.log.Error(err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		e.log.Error(err.Error())
		return 0, err
	}

	blockHex := result["result"].(string)
	blockNumber, err := strconv.ParseInt(blockHex[2:], 16, 64)
	if err != nil {
		e.log.Error(err.Error())
		return 0, err
	}

	return int(blockNumber), nil
}

func (e *Ethereum) GetBlockTransactions(ctx context.Context, blockNumber int) ([]parser.Transaction, error) {
	e.log.Debug(fmt.Sprintf("Executing GetBlockTransactions. Block: %v", blockNumber))

	payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x%x", true],"id":1}`, blockNumber)

	resp, err := http.Post(e.cliUrl, "application/json", strings.NewReader(payload))
	if err != nil {
		e.log.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		e.log.Error(err.Error())
		return nil, err
	}

	block := result["result"].(map[string]interface{})
	txs := block["transactions"].([]interface{})

	var transactions []parser.Transaction
	for _, tx := range txs {
		txMap := tx.(map[string]interface{})

		hash, ok := txMap["hash"].(string)
		if !ok {
			e.log.Error(fmt.Sprintf("%s", txMap["hash"]))
			continue
		}

		from, ok := txMap["from"].(string)
		if !ok {
			e.log.Error(fmt.Sprintf("%s", txMap["from"]))
			from = ""
		}

		to, ok := txMap["to"].(string)
		if !ok {
			e.log.Error(fmt.Sprintf("%s", txMap["to"]))
			to = ""
		}

		value, ok := txMap["value"].(string)
		if !ok {
			e.log.Error(fmt.Sprintf("%s", txMap["value"]))
			value = "0x0"
		}

		transactions = append(transactions, parser.Transaction{
			Hash:        hash,
			From:        from,
			To:          to,
			Value:       value,
			BlockNumber: blockNumber,
		})
	}

	return transactions, nil
}
