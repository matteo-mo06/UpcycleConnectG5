package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"upcycle_connect-api/internal/db"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roles, err := db.GetAllRoles()
	if err != nil {
		fmt.Println("GetRoles error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch roles"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(roles)
}

func GetUserRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	roles, err := db.GetUserRoles(id)
	if err != nil {
		fmt.Println("GetUserRoles error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch user roles"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(roles)
}

func AddRoleToUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	var body struct {
		RoleID string `json:"role_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.RoleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "role_id is required"})
		return
	}

	_, err = db.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("AddRoleToUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find user"})
		return
	}

	err = db.AddUserRole(id, body.RoleID)
	if err != nil {
		fmt.Println("AddRoleToUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to add role"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "role added successfully"})
}

func RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.PathValue("id")
	roleID := r.PathValue("role_id")

	if userID == "" || roleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id or role id"})
		return
	}

	deleted, err := db.RemoveUserRole(userID, roleID)
	if err != nil {
		fmt.Println("RemoveRoleFromUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to remove role"})
		return
	}
	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "role not found for this user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "role removed successfully"})
}
