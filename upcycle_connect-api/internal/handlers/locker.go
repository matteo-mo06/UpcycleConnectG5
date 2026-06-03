package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/models"
)

var accessCodeRe = regexp.MustCompile(`^\d{8}$`)

func GetLockers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := db.GetAllLockers(r.URL.Query().Get("status"))
	if err != nil {
		fmt.Println("GetLockers error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch lockers"})
		return
	}
	if list == nil {
		list = []models.Locker{}
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func CreateLocker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	locker, err := db.CreateLocker()
	if err != nil {
		fmt.Println("CreateLocker error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create locker"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(locker)
}

func UpdateLockerAccessCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	var body struct {
		AccessCode string `json:"access_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || !accessCodeRe.MatchString(body.AccessCode) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le code d'accès doit contenir exactement 8 chiffres"})
		return
	}
	if err := db.UpdateLockerAccessCode(id, body.AccessCode); err != nil {
		if err.Error() == "locker libre" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "impossible de modifier le code d'un casier libre"})
			return
		}
		fmt.Println("UpdateLockerAccessCode error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update access code"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "code mis à jour"})
}

func DeleteLocker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	if err := db.DeleteLocker(id); err != nil {
		if err.Error() == "locker occupé" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "ce casier est actuellement occupé"})
			return
		}
		fmt.Println("DeleteLocker error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete locker"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "casier supprimé"})
}
