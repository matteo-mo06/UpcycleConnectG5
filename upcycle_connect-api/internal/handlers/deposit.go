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
	"upcycle_connect-api/internal/utils"
)

func GetMyDepositRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	list, err := db.GetUserDepositRequests(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(list)
}

func GetPendingDepositRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list, err := db.GetPendingDepositRequests()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
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
	if ann.StateAnnouncement != "Vendu" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "seules les annonces vendues peuvent faire l'objet d'une demande de dépôt"})
		return
	}
	if ann.Request == 1 {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "une demande de dépôt est déjà en cours"})
		return
	}

	if _, err := db.GetAvailableLocker(); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Aucun casier disponible pour le moment. Veuillez réessayer ultérieurement."})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	if err := db.SetDepositRequest(announcementID, 1); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création de la demande"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "demande de dépôt créée, en attente de validation"})
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
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if err := db.SetDepositRequest(announcementID, 0); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ValidateDepositRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	announcementID := r.PathValue("id")

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
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "aucun casier disponible"})
			return
		}
		fmt.Println("ReserveLocker error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de l'assignation du casier"})
		return
	}

	if err := db.SetDepositRequest(announcementID, 0); err != nil {
		_ = db.UnassignLocker(announcementID)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if ownerID, err := db.GetAnnouncementOwnerID(announcementID); err == nil {
		if err := db.AwardScore(ownerID, "deposit_validated", announcementID); err != nil {
			fmt.Println("AwardScore deposit_validated error:", err)
		}
		go utils.SendPushNotification(db.GetOnesignalPlayerID(ownerID), "Dépôt validé", fmt.Sprintf("Votre demande de dépôt a été validée. Déposez votre objet dans le casier n°%d.", locker.Number))
	}

	if buyerID, err := db.GetAnnouncementBuyerID(announcementID); err == nil && buyerID != "" {
		go utils.SendPushNotification(db.GetOnesignalPlayerID(buyerID), "Casier assigné", fmt.Sprintf("L'objet que vous avez acheté sera déposé dans le casier n°%d. Code d'accès : %s", locker.Number, accessCode))
	}

	w.WriteHeader(http.StatusNoContent)
}

func RejectDepositRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	announcementID := r.PathValue("id")

	if err := db.SetDepositRequest(announcementID, 0); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if ownerID, err := db.GetAnnouncementOwnerID(announcementID); err == nil {
		go utils.SendPushNotification(db.GetOnesignalPlayerID(ownerID), "Dépôt refusé", "Votre demande de dépôt n'a pas été acceptée.")
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
