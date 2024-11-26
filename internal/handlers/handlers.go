package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmsilvadev/tx-parser/pkg/parser"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type handler struct {
	parser parser.Parser
}

func New(p parser.Parser) *handler {
	return &handler{parser: p}
}

func (h *handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status: "success",
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func (h *handler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := h.parser.GetCurrentBlock()
	response := Response{
		Status: "success",
		Data:   map[string]int{"currentBlock": block},
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func (h *handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		response := Response{
			Status:  "error",
			Message: "address is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	success := h.parser.Subscribe(address)
	response := Response{
		Status: "error",
		Data:   map[string]bool{"subscribed": success},
	}

	if !success {
		response.Message = "Address already subscribed or invalid"
	}
	writeJSONResponse(w, http.StatusBadRequest, response)
}

func (h *handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		response := Response{
			Status:  "error",
			Message: "address is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	transactions := h.parser.GetTransactions(address)
	if transactions == nil {
		transactions = []parser.Transaction{}
	}

	response := Response{
		Status: "success",
		Data:   transactions,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
