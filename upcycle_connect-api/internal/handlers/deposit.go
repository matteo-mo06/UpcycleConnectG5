package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
)

func GetMyDepositRequests(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	list, err := db.GetUserDepositRequests(userID)
	if err != nil {
		http.Error(w, "GetUserDepositRequests error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func GetPendingDepositRequests(w http.ResponseWriter, r *http.Request) {
	list, err := db.GetPendingDepositRequests()
	if err != nil {
		http.Error(w, "GetPendingDepositRequests error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func CreateDepositRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	announcementID := r.PathValue("id")

	owner, err := db.IsAnnouncementOwner(announcementID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if !owner {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
		return
	}

	ann, err := db.GetAnnouncementById(announcementID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "annonce introuvable"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if ann.State_annoucement != "Vendu" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "seules les annonces vendues peuvent faire l'objet d'une demande de dépôt"})
		return
	}

	assigned, err := db.IsLockerAssigned(announcementID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if assigned {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "un casier est déjà assigné à cette annonce"})
		return
	}

	accessCode, err := generateAccessCode()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate access code"})
		return
	}

	locker, err := db.ReserveLocker(announcementID, accessCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Aucun casier disponible pour le moment. Veuillez réessayer ultérieurement.",
			})
			return
		}
		fmt.Println("ReserveLocker error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la réservation du casier"})
		return
	}

	if err := db.SetDepositRequest(announcementID, 1); err != nil {
		_ = db.UnassignLocker(announcementID)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création de la demande"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message":       "casier réservé",
		"locker_number": locker.Number,
	})
}

func CancelDepositRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	announcementID := r.PathValue("id")

	owner, err := db.IsAnnouncementOwner(announcementID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if !owner {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
		return
	}

	if err := db.UnassignLocker(announcementID); err != nil {
		fmt.Println("CancelDepositRequest UnassignLocker error:", err)
	}
	if err := db.SetDepositRequest(announcementID, 0); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ValidateDepositRequest(w http.ResponseWriter, r *http.Request) {
	announcementID := r.PathValue("id")

	// Locker is already assigned at request creation time — just confirm deposit.
	if err := db.SetDepositRequest(announcementID, 0); err != nil {
		http.Error(w, "ValidateDepositRequest error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if ownerID, err := db.GetAnnouncementOwnerID(announcementID); err == nil {
		if err := db.AwardScore(ownerID, "deposit_validated", announcementID); err != nil {
			fmt.Println("AwardScore deposit_validated error:", err)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func RejectDepositRequest(w http.ResponseWriter, r *http.Request) {
	announcementID := r.PathValue("id")

	if err := db.UnassignLocker(announcementID); err != nil {
		fmt.Println("RejectDepositRequest UnassignLocker error:", err)
	}
	if err := db.SetDepositRequest(announcementID, 0); err != nil {
		http.Error(w, "RejectDepositRequest error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func generateAccessCode() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	code := make([]byte, 8)
	for i := range code {
		code[i] = '0' + b[i]%10
	}
	return string(code), nil
}
