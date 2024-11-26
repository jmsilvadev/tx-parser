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
	if r.Method != http.MethodGet {
		response := Response{
			Status:  "error",
			Message: "method not allowed",
		}
		writeJSONResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	block := h.parser.GetCurrentBlock()
	response := Response{
		Status: "success",
		Data:   map[string]int{"currentBlock": block},
	}
	writeJSONResponse(w, http.StatusOK, response)
}

func (h *handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := Response{
			Status:  "error",
			Message: "method not allowed",
		}
		writeJSONResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	var reqBody struct {
		Address string `json:"address"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil || reqBody.Address == "" {
		response := Response{
			Status:  "error",
			Message: "address is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	success := h.parser.Subscribe(reqBody.Address)
	response := Response{
		Status: "success",
		Data:   map[string]bool{"subscribed": success},
	}

	if !success {
		response.Status = "error"
		response.Message = "Address already subscribed or invalid"
	}
	writeJSONResponse(w, http.StatusBadRequest, response)
}

func (h *handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := Response{
			Status:  "error",
			Message: "method not allowed",
		}
		writeJSONResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

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

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "error",
		Message: "route not found",
	}
	writeJSONResponse(w, http.StatusNotFound, response)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
