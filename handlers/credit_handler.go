package handlers

import (
	"encoding/json"
	"net/http"

	"bank-api/models"
	"bank-api/services"
	"github.com/gorilla/mux"
	"strconv"
)


type ScheduleService interface {
    GetSchedule(creditID int) ([]*models.PaymentSchedule, error)
}

// CreditHandler содержит зависимости для работы с кредитами.
type CreditHandler struct {
	creditService services.CreditService
	scheduleService ScheduleService
}

// NewCreditHandler возвращает новый экземпляр CreditHandler.
func NewCreditHandler(creditService services.CreditService) *CreditHandler {
	return &CreditHandler{creditService: creditService}
}

// ApplyForCredit обрабатывает POST-запрос на оформление кредита.
// URL: /credits
func (h *CreditHandler) ApplyForCredit(w http.ResponseWriter, r *http.Request) {
	var credit models.Credit
	if err := json.NewDecoder(r.Body).Decode(&credit); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.creditService.ApplyForCredit(&credit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credit)
}

// GET /credits/{creditId}/schedule
func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["creditId"])
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}
	schedule, err := h.scheduleService.GetSchedule(id)
	if err != nil {
		http.Error(w, "Error fetching schedule: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

