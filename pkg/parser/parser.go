package parser

import "context"

type Parser interface {
	// last parsed block
	GetCurrentBlock(context.Context) int
	// add address to observer
	Subscribe(context.Context, string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(context.Context, string) []Transaction
	// routine to fetch the transactions each 12 seconds
	UpdateBlockNumber(context.Context)
}

type Transaction struct {
	Hash        string `json:"hash,omitempty"`
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	Value       string `json:"value,omitempty"`
	BlockNumber int    `json:"block_number,omitempty"`
}
