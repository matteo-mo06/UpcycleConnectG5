package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/sse"
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
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	announcementID := r.PathValue("id")

	owner, err := db.IsAnnouncementOwner(announcementID, userID)
	if err != nil {
		http.Error(w, "IsAnnouncementOwner error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !owner {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := db.SetDepositRequest(announcementID, 1); err != nil {
		http.Error(w, "SetDepositRequest error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func CancelDepositRequest(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	announcementID := r.PathValue("id")

	owner, err := db.IsAnnouncementOwner(announcementID, userID)
	if err != nil {
		http.Error(w, "IsAnnouncementOwner error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !owner {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := db.SetDepositRequest(announcementID, 0); err != nil {
		http.Error(w, "SetDepositRequest error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ValidateDepositRequest(w http.ResponseWriter, r *http.Request) {
	announcementID := r.PathValue("id")

	locker, err := db.GetAvailableLocker()
	if err != nil {
		http.Error(w, "no locker available", http.StatusConflict)
		return
	}

	accessCode, err := generateAccessCode()
	if err != nil {
		http.Error(w, "failed to generate access code", http.StatusInternalServerError)
		return
	}

	if err := db.AssignLocker(announcementID, locker.ID, accessCode); err != nil {
		http.Error(w, "AssignLocker error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if ownerID, err := db.GetAnnouncementOwnerID(announcementID); err == nil {
		sse.Default.Notify(ownerID, fmt.Sprintf(
			`{"type":"deposit_validated","locker_number":%d}`,
			locker.Number,
		))
	}

	w.WriteHeader(http.StatusNoContent)
}

func RejectDepositRequest(w http.ResponseWriter, r *http.Request) {
	announcementID := r.PathValue("id")

	if err := db.SetDepositRequest(announcementID, 0); err != nil {
		http.Error(w, "RejectDepositRequest error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if ownerID, err := db.GetAnnouncementOwnerID(announcementID); err == nil {
		sse.Default.Notify(ownerID, `{"type":"deposit_rejected"}`)
	}

	w.WriteHeader(http.StatusNoContent)
}

func generateAccessCode() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%08X", b), nil
}
