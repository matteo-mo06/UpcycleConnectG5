package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
)

func SubmitProfessionalRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	_, err := db.GetPendingRequestByUser(userID)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "une demande est déjà en attente"})
		return
	}
	if err != sql.ErrNoRows {
		fmt.Println("SubmitProfessionalRequest error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to check existing request"})
		return
	}

	req := models.ProfessionalRequest{
		Id_request: uuid.New().String(),
		Id_user:    userID,
	}

	if err := db.CreateProfessionalRequest(req); err != nil {
		fmt.Println("SubmitProfessionalRequest error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create request"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "demande envoyée avec succès"})
}

func GetProfessionalRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := r.URL.Query().Get("status")

	requests, err := db.GetProfessionalRequests(status)
	if err != nil {
		fmt.Println("GetProfessionalRequests error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch requests"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(requests)
}

func ValidateProfessionalRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing request id"})
		return
	}

	req, err := db.GetProfessionalRequestById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "request not found"})
			return
		}
		fmt.Println("ValidateProfessionalRequest error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get request"})
		return
	}

	if req.Status != "pending" {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette demande a déjà été traitée"})
		return
	}

	if err := db.UpdateProfessionalRequestStatus(id, "approved"); err != nil {
		fmt.Println("ValidateProfessionalRequest error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update request status"})
		return
	}

	artisanRole, err := db.GetRoleByName(config.ArtisanRoleName)
	if err == nil {
		_ = db.AddUserRole(req.Id_user, artisanRole.Id_role)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "demande approuvée, rôle artisan attribué"})
}

func RejectProfessionalRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing request id"})
		return
	}

	req, err := db.GetProfessionalRequestById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "request not found"})
			return
		}
		fmt.Println("RejectProfessionalRequest error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get request"})
		return
	}

	if req.Status != "pending" {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette demande a déjà été traitée"})
		return
	}

	if err := db.UpdateProfessionalRequestStatus(id, "rejected"); err != nil {
		fmt.Println("RejectProfessionalRequest error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update request status"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "demande rejetée"})
}
