package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/models"
)

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	perms, err := db.GetAllPermissions()
	if err != nil {
		fmt.Println("GetPermissions error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch permissions"})
		return
	}

	if perms == nil {
		perms = []models.Permission{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(perms)
}

func GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roleID := r.PathValue("id")
	if roleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing role id"})
		return
	}

	perms, err := db.GetPermissionsByRole(roleID)
	if err != nil {
		fmt.Println("GetRolePermissions error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch role permissions"})
		return
	}

	if perms == nil {
		perms = []models.Permission{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(perms)
}

func AddRolePermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roleID := r.PathValue("id")
	if roleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing role id"})
		return
	}

	var body struct {
		PermissionId string `json:"permission_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.PermissionId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "permission_id requis"})
		return
	}

	if err := db.AddRolePermission(roleID, body.PermissionId); err != nil {
		fmt.Println("AddRolePermission error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to add permission"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "permission ajoutée au rôle"})
}

func RemoveRolePermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roleID := r.PathValue("id")
	permissionID := r.PathValue("permission_id")
	if roleID == "" || permissionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing role id or permission id"})
		return
	}

	if err := db.RemoveRolePermission(roleID, permissionID); err != nil {
		fmt.Println("RemoveRolePermission error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to remove permission"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "permission retirée du rôle"})
}
