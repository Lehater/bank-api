package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/services"

	"github.com/gorilla/mux"
)

// AnalyticsHandler содержит зависимости для аналитики.
type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
}

// NewAnalyticsHandler создаёт новый обработчик аналитики.
func NewAnalyticsHandler(as services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: as,
	}
}

// GetAnalytics обрабатывает GET /analytics и возвращает данные по доходам/расходам за текущий месяц.
// Ожидается, что идентификатор пользователя (userID) уже присутствует в контексте (например, через AuthMiddleware).
func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok || userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	data, err := h.analyticsService.GetAnalytics(userID)
	if err != nil {
		http.Error(w, "Error retrieving analytics: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// PredictBalance обрабатывает GET /accounts/{accountId}/predict?days=N и возвращает прогноз баланса.
func (h *AnalyticsHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountIDStr := vars["accountId"]
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	daysStr := r.URL.Query().Get("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		http.Error(w, "Invalid days parameter", http.StatusBadRequest)
		return
	}

	prediction, err := h.analyticsService.PredictBalance(accountID, days)
	if err != nil {
		http.Error(w, "Error predicting balance: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]float64{"predicted_balance": prediction}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
