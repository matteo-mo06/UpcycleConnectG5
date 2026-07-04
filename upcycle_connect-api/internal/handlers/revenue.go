package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/models"
)

func GetRevenueSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	summary, err := db.GetRevenueSummary()
	if err != nil {
		fmt.Println("GetRevenueSummary error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la récupération des revenus"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(summary)
}

func GetRevenueTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	transactions, total, err := db.GetRevenueTransactions(page, limit)
	if err != nil {
		fmt.Println("GetRevenueTransactions error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la récupération des transactions"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

func GetAdPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	list, total, err := db.GetAdPaymentsPaginated(page, limit)
	if err != nil {
		fmt.Println("GetAdPayments error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if list == nil {
		list = []models.AdPayment{}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"payments": list,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func GetFormationPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	list, total, err := db.GetFormationPaymentsPaginated(page, limit)
	if err != nil {
		fmt.Println("GetFormationPayments error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if list == nil {
		list = []models.FormationPayment{}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"payments": list,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func GetCommissionRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rate, err := db.GetCommissionRate()
	if err != nil {
		fmt.Println("GetCommissionRate error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la récupération du taux"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]float64{"commission_rate": rate})
}

func UpdateCommissionRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Rate float64 `json:"commission_rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}
	if body.Rate < 0 || body.Rate > 100 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le taux doit être entre 0 et 100"})
		return
	}

	if err := db.SetCommissionRate(body.Rate); err != nil {
		fmt.Println("UpdateCommissionRate error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la mise à jour du taux"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]float64{"commission_rate": body.Rate})
}
