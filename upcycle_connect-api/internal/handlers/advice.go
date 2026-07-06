package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
)

func GetAdvices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, limit, offset := parsePage(r, 15)

	advices, total, err := db.GetAdvices(
		r.URL.Query().Get("search"),
		r.URL.Query().Get("id_category"),
		limit, offset,
	)
	if err != nil {
		fmt.Println("GetAdvices error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch advices"})
		return
	}

	if advices == nil {
		advices = []models.Advice{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(advices, total, page, limit))
}

func GetAdvice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	advice, err := db.GetAdviceById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "conseil introuvable"})
			return
		}
		fmt.Println("GetAdvice error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get advice"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(advice)
}

func GetAdviceDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	docs, err := db.GetDocumentsByCategory(id)
	if err != nil {
		fmt.Println("GetAdviceDocuments error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get documents"})
		return
	}
	if docs == nil {
		docs = []models.Document{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(docs)
}

func CreateAdvice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var req models.CreateAdviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le titre est requis"})
		return
	}

	if req.Priority < 1 || req.Priority > 3 {
		req.Priority = 2
	}
	if req.DateAdvice != nil && *req.DateAdvice == "" {
		req.DateAdvice = nil
	}

	if err := db.CreateAdvice(userID, req.Title, req.Description, req.IdCategory, req.DateAdvice, req.Priority, req.PhotoURLs); err != nil {
		fmt.Println("CreateAdvice error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create advice"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "conseil créé"})
}

func UpdateMyAdvice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	id := r.PathValue("id")

	creatorID, err := db.GetAdviceCreatorID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "conseil introuvable"})
			return
		}
		fmt.Println("UpdateMyAdvice fetch error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find advice"})
		return
	}

	isAdmin := slices.Contains(roles, config.RoleAdmin)

	if creatorID != userID && !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous n'êtes pas créateur de ce conseil"})
		return
	}

	var req models.UpdateAdviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le titre est requis"})
		return
	}

	if req.Priority < 1 || req.Priority > 3 {
		req.Priority = 2
	}
	if req.DateAdvice != nil && *req.DateAdvice == "" {
		req.DateAdvice = nil
	}

	if err := db.UpdateAdvice(id, req.Title, req.Description, req.IdCategory, req.DateAdvice, req.Priority); err != nil {
		fmt.Println("UpdateMyAdvice error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update advice"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "conseil mis à jour"})
}

func DeleteMyAdvice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	perms, _ := r.Context().Value(middleware.ContextPermissions).([]string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	id := r.PathValue("id")

	creatorID, err := db.GetAdviceCreatorID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "conseil introuvable"})
			return
		}
		fmt.Println("DeleteMyAdvice fetch error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find advice"})
		return
	}

	isOwner := creatorID == userID
	canManage := slices.Contains(perms, "manage_advice")
	isAdmin := slices.Contains(roles, config.RoleAdmin)

	if !isOwner && !canManage && !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès non autorisé"})
		return
	}

	if err := db.DeleteAdvice(id); err != nil {
		fmt.Println("DeleteMyAdvice error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete advice"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "conseil supprimé"})
}

func GetDeletedAdvices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	advices, err := db.GetDeletedAdvices()
	if err != nil {
		fmt.Println("GetDeletedAdvices error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch advices"})
		return
	}

	if advices == nil {
		advices = []models.DeletedAdvice{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(advices)
}
