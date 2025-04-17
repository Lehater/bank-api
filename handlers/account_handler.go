package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/models"
	"bank-api/services"
	"github.com/gorilla/mux"
)

// AccountHandler содержит зависимости для работы со счетами.
type AccountHandler struct {
	accountService services.AccountService
}

// NewAccountHandler возвращает новый экземпляр AccountHandler.
func NewAccountHandler(accountService services.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// CreateAccount обрабатывает POST-запрос на создание нового счета.
// URL: /accounts
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.accountService.CreateAccount(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// Deposit обрабатывает POST-запрос для пополнения счета.
// URL: /accounts/{id}/deposit
func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор счета из URL
	vars := mux.Vars(r)
	accountID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный идентификатор счета", http.StatusBadRequest)
		return
	}

	var payload struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.accountService.Deposit(accountID, payload.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Если операция выполнена успешно, можно вернуть статус 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}

// POST /transfer
func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromAccountID int     `json:"from_account_id"`
		ToAccountID   int     `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := h.accountService.Transfer(req.FromAccountID, req.ToAccountID, req.Amount); err != nil {
		http.Error(w, "Transfer failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}