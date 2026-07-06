package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/google/uuid"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
)

func GetFormationSteps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	if _, err := db.GetFormationById(id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation introuvable"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	steps, err := db.GetFormationSteps(id)
	if err != nil {
		fmt.Println("GetFormationSteps error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	if steps == nil {
		steps = []models.FormationStep{}
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(steps)
}

func CreateFormationStep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	ownerID, err := db.GetFormationOwnerID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation introuvable"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	if ownerID != userID && !slices.Contains(roles, config.RoleAdmin) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "réservé au créateur de la formation"})
		return
	}

	var req models.CreateFormationStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps invalide"})
		return
	}
	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le titre est requis"})
		return
	}

	s := models.FormationStep{
		IdStep:      uuid.New().String(),
		IdFormation: id,
		Title:       req.Title,
		Description: req.Description,
		Status:      "todo",
		StepOrder:   req.StepOrder,
	}
	if err := db.CreateFormationStep(s); err != nil {
		fmt.Println("CreateFormationStep error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(s)
}

func UpdateFormationStep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")
	stepID := r.PathValue("step_id")

	ownerID, err := db.GetFormationOwnerID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation introuvable"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	if ownerID != userID && !slices.Contains(roles, config.RoleAdmin) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "réservé au créateur de la formation"})
		return
	}

	var req models.UpdateFormationStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps invalide"})
		return
	}
	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le titre est requis"})
		return
	}

	s := models.FormationStep{
		IdStep:      stepID,
		Title:       req.Title,
		Description: req.Description,
		StepOrder:   req.StepOrder,
	}
	if err := db.UpdateFormationStep(s); err != nil {
		fmt.Println("UpdateFormationStep error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la mise à jour"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "étape mise à jour"})
}

func UpdateFormationStepStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")
	stepID := r.PathValue("step_id")

	ownerID, err := db.GetFormationOwnerID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation introuvable"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	if ownerID != userID && !slices.Contains(roles, config.RoleAdmin) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "réservé au créateur de la formation"})
		return
	}

	var req models.UpdateFormationStepStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps invalide"})
		return
	}
	if !slices.Contains([]string{"todo", "in_progress", "done"}, req.Status) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "statut invalide (todo, in_progress, done)"})
		return
	}

	if err := db.UpdateFormationStepStatus(stepID, req.Status); err != nil {
		fmt.Println("UpdateFormationStepStatus error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la mise à jour"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "statut mis à jour"})
}

func DeleteFormationStep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")
	stepID := r.PathValue("step_id")

	ownerID, err := db.GetFormationOwnerID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation introuvable"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	if ownerID != userID && !slices.Contains(roles, config.RoleAdmin) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "réservé au créateur de la formation"})
		return
	}

	if err := db.DeleteFormationStep(stepID); err != nil {
		fmt.Println("DeleteFormationStep error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la suppression"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "étape supprimée"})
}
