package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
	"upcycle_connect-api/internal/utils"
)

func GetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	user, err := db.GetUserById(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("GetMe error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get user"})
		return
	}

	roles, _ := db.GetUserRoleNames(userID)
	if roles == nil {
		roles = []string{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user.ToLoginResponse(roles, []string{}))
}

func GetMyStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	stats, err := db.GetUserStats(userID)
	if err != nil {
		fmt.Println("GetMyStats error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get stats"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}

func UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	allowed := map[string]bool{"active": true, "suspended": true, "blacklisted": true}
	if !allowed[body.Status] {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "status invalide (active | suspended | blacklisted)"})
		return
	}

	if err := db.UpdateUserStatus(id, body.Status); err != nil {
		fmt.Println("UpdateUserStatus error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update status"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "statut mis à jour"})
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	user, err := db.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("GetUserById error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user.ToResponse())
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	users, err := db.GetAllUsersWithRoles(firstName, lastName)
	if err != nil {
		fmt.Println("GetUsers error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch users"})
		return
	}

	if users == nil {
		users = []models.UserListItem{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	_, err := db.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("DeleteUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find user"})
		return
	}

	err = db.DeleteUser(id)
	if err != nil {
		fmt.Println("DeleteUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "user deleted successfully"})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	_, err := db.GetUserById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		fmt.Println("UpdateUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find user"})
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	user.Id_user = id

	err = db.UpdateUser(user)
	if err != nil {
		fmt.Println("UpdateUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "user updated successfully", "user": user.ToResponse()})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	user.Email = strings.TrimSpace(user.Email)
	user.First_name = strings.TrimSpace(user.First_name)
	user.Last_name = strings.TrimSpace(user.Last_name)

	err = utils.ValidateUserCreation(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password_user)
	if err != nil {
		fmt.Println("CreateUser hash error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to hash password"})
		return
	}

	user.Id_user = uuid.New().String()
	user.Password_user = hashedPassword

	err = db.CreateUser(user)
	if err != nil {
		fmt.Println("CreateUser error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create user"})
		return
	}

	role, err := db.GetRoleByName("user")
	if err != nil {
		fmt.Println("CreateUser GetRoleByName error:", err)
	} else {
		err = db.AddUserRole(user.Id_user, role.Id_role)
		if err != nil {
			fmt.Println("CreateUser AddUserRole error:", err)
		}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "user created successfully",
		"user":    user.ToResponse(),
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	payload.Email = strings.TrimSpace(payload.Email)

	if payload.Email == "" || payload.Password_user == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "email and password are required"})
		return
	}

	user, err := db.GetUserByEmail(payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
			return
		}
		fmt.Println("Login error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch user"})
		return
	}

	err = utils.CheckPassword(payload.Password_user, user.Password_user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	if user.Status == "blacklisted" {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "compte banni"})
		return
	}

	if user.Status == "suspended" {
		expired, err := db.IsSuspensionExpired(user.Id_user)
		if err != nil {
			fmt.Println("Login IsSuspensionExpired error:", err)
		}
		if expired {
			_ = db.UpdateUserStatus(user.Id_user, "active")
		} else {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "compte suspendu"})
			return
		}
	}

	roles, _ := db.GetUserRoleNames(user.Id_user)
	if roles == nil {
		roles = []string{}
	}

	permissions, _ := db.GetUserPermissions(user.Id_user)
	if permissions == nil {
		permissions = []string{}
	}

	token, err := utils.GenerateToken(user.Id_user, roles, permissions)
	if err != nil {
		fmt.Println("Login GenerateToken error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "login successful",
		"token":   token,
		"user":    user.ToLoginResponse(roles, permissions),
	})
}
