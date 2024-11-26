package parser

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
	// routine to fetch the transactions each 12 seconds
	UpdateBlockNumber()
}

type Transaction struct {
	Hash        string `json:"hash,omitempty"`
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	Value       string `json:"value,omitempty"`
	BlockNumber int    `json:"block_number,omitempty"`
}
