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
)

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	id := r.PathValue("id")

	doc, err := db.GetDocumentById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "document introuvable"})
			return
		}
		fmt.Println("DeleteDocument fetch error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find document"})
		return
	}

	isOwner := doc.IdUser == userID
	isAdmin := slices.Contains(roles, config.RoleAdmin)

	if !isOwner && !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès non autorisé"})
		return
	}

	if err := db.DeleteDocument(id); err != nil {
		fmt.Println("DeleteDocument error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete document"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "document supprimé"})
}
