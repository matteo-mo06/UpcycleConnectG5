package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
	"upcycle_connect-api/internal/utils"
)

func GetFormations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	page, limit, offset := parsePage(r, 15)

	formations, total, err := db.GetFormations(userID, r.URL.Query().Get("search"), r.URL.Query().Get("level"), r.URL.Query().Get("id_category"), limit, offset)
	if err != nil {
		fmt.Println("GetFormations error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch formations"})
		return
	}

	if formations == nil {
		formations = []models.Formation{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(formations, total, page, limit))
}

func GetFormationById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	perms, _ := r.Context().Value(middleware.ContextPermissions).([]string)

	formation, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		fmt.Println("GetFormationById error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get formation"})
		return
	}

	canSee := formation.Status == "approved" ||
		(formation.IdCreator != nil && *formation.IdCreator == userID) ||
		slices.Contains(perms, "manage_formations")

	if !canSee {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès non autorisé"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(formation)
}

func CreateFormation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var req models.CreateFormationRequest
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

	if req.Date == nil || *req.Date == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "la date est requise"})
		return
	}

	if req.Location == nil || *req.Location == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le lieu est requis"})
		return
	}

	if req.Date != nil && *req.Date != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05", *req.Date)
		if err != nil {
			parsed, err = time.Parse("2006-01-02T15:04", *req.Date)
		}
		if err != nil || parsed.Before(time.Now()) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "la date doit être dans le futur"})
			return
		}
	}

	validLevels := []string{"beginner", "intermediate", "advanced"}
	if req.Level == "" {
		req.Level = "beginner"
	}
	if !slices.Contains(validLevels, req.Level) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "niveau invalide (beginner, intermediate, advanced)"})
		return
	}

	f := models.Formation{
		IdFormation:          uuid.New().String(),
		TitleFormation:       req.Title,
		DescriptionFormation: req.Description,
		DateFormation:        req.Date,
		LocationFormation:    req.Location,
		Capacity:             req.Capacity,
		Level:                req.Level,
		DurationHours:        req.DurationH,
		IdCategory:           req.IdCategory,
		Price:                req.Price,
		IdCreator:            &userID,
		IdFormateur:          &userID,
	}

	if err := db.CreateFormation(f); err != nil {
		fmt.Println("CreateFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create formation"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "formation créée, en attente de validation", "formation": f})
}

func UpdateMyFormation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	formation, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find formation"})
		return
	}

	if formation.IdCreator == nil || *formation.IdCreator != userID {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous n'êtes pas créateur de cette formation"})
		return
	}

	if formation.Status == "approved" {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "une formation approuvée ne peut pas être modifiée directement"})
		return
	}

	var req models.UpdateFormationRequest
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

	if req.Date == nil || *req.Date == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "la date est requise"})
		return
	}

	if req.Location == nil || *req.Location == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le lieu est requis"})
		return
	}

	if req.Date != nil && *req.Date != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05", *req.Date)
		if err != nil {
			parsed, err = time.Parse("2006-01-02T15:04", *req.Date)
		}
		if err != nil || parsed.Before(time.Now()) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "la date doit être dans le futur"})
			return
		}
	}

	formation.TitleFormation = req.Title
	formation.DescriptionFormation = req.Description
	formation.DateFormation = req.Date
	formation.LocationFormation = req.Location
	formation.Capacity = req.Capacity
	formation.Level = req.Level
	formation.DurationHours = req.DurationH
	formation.IdCategory = req.IdCategory
	formation.Price = req.Price

	if err := db.UpdateFormation(formation); err != nil {
		fmt.Println("UpdateMyFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update formation"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "formation mise à jour, resoumise pour validation"})
}

func DeleteMyFormation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	perms, _ := r.Context().Value(middleware.ContextPermissions).([]string)
	id := r.PathValue("id")

	formation, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find formation"})
		return
	}

	isOwner := formation.IdCreator != nil && *formation.IdCreator == userID
	canManage := slices.Contains(perms, "manage_formations")

	if !isOwner && !canManage {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès non autorisé"})
		return
	}

	if err := db.DeleteFormation(id); err != nil {
		fmt.Println("DeleteMyFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete formation"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "formation supprimée"})
}

func GetPendingFormations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, limit, offset := parsePage(r, 20)
	formations, total, err := db.GetAllFormations("", "pending", "", limit, offset)
	if err != nil {
		fmt.Println("GetPendingFormations error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch pending formations"})
		return
	}

	if formations == nil {
		formations = []models.Formation{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(formations, total, page, limit))
}

func ApproveFormation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	_, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find formation"})
		return
	}

	if err := db.ApproveFormation(id); err != nil {
		fmt.Println("ApproveFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to approve formation"})
		return
	}

	go utils.SendPushNotification(db.GetOnesignalPlayerID(db.GetFormationCreatorID(id)), "Formation approuvée", "Votre formation a été approuvée et est maintenant visible.")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "formation approuvée"})
}

func RejectFormation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	_, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find formation"})
		return
	}

	var req models.RejectFormationRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := db.RejectFormation(id, req.Reason); err != nil {
		fmt.Println("RejectFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to reject formation"})
		return
	}

	go utils.SendPushNotification(db.GetOnesignalPlayerID(db.GetFormationCreatorID(id)), "Formation refusée", "Votre formation n'a pas été approuvée.")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "formation rejetée"})
}

func RegisterForFormation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	formation, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find formation"})
		return
	}

	if formation.Status != "approved" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette formation n'est pas encore disponible"})
		return
	}

	if formation.Capacity != nil && *formation.Capacity > 0 && formation.InscriptionCount >= *formation.Capacity {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "la formation est complète"})
		return
	}

	if formation.Price != nil && *formation.Price > 0 {
		w.WriteHeader(http.StatusPaymentRequired)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "paiement requis", "price": *formation.Price})
		return
	}

	if err := db.RegisterUserForFormation(userID, id); err != nil {
		fmt.Println("RegisterForFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to register for formation"})
		return
	}

	if err := db.AwardScore(userID, "formation_registration", id); err != nil {
		fmt.Println("AwardScore formation_registration error:", err)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "inscription réussie"})
}

func GetFormationParticipants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)

	formation, err := db.GetFormationById(id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation introuvable"})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	isCreator := formation.IdCreator != nil && *formation.IdCreator == userID
	isAdmin := slices.Contains(roles, config.RoleAdmin)
	if !isCreator && !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}

	participants, err := db.GetFormationParticipants(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if participants == nil {
		participants = []models.Participant{}
	}
	_ = json.NewEncoder(w).Encode(participants)
}

func UnregisterFromFormation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	formation, err := db.GetFormationById(id)
	if err == nil && formation.Price != nil && *formation.Price > 0 {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "impossible de se désinscrire d'une formation payante"})
		return
	}

	if err := db.UnregisterUserFromFormation(userID, id); err != nil {
		fmt.Println("UnregisterFromFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to unregister from formation"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "désinscription réussie"})
}

func GetUserFormations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	formations, err := db.GetUserRegisteredFormations(userID)
	if err != nil {
		fmt.Println("GetUserFormations error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch user formations"})
		return
	}

	if formations == nil {
		formations = []models.Formation{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(formations)
}

func GetMyCreatedFormations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	page, limit, offset := parsePage(r, 15)

	formations, total, err := db.GetMyCreatedFormationsPaginated(
		userID,
		r.URL.Query().Get("search"),
		r.URL.Query().Get("status"),
		r.URL.Query().Get("level"),
		r.URL.Query().Get("id_category"),
		limit,
		offset,
	)
	if err != nil {
		fmt.Println("GetMyCreatedFormations error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch created formations"})
		return
	}

	if formations == nil {
		formations = []models.Formation{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(formations, total, page, limit))
}


func GetAllFormationsAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, limit, offset := parsePage(r, 20)
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	level := r.URL.Query().Get("level")

	formations, total, err := db.GetAllFormations(search, status, level, limit, offset)
	if err != nil {
		fmt.Println("GetAllFormationsAdmin error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch formations"})
		return
	}

	if formations == nil {
		formations = []models.Formation{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(formations, total, page, limit))
}

func GetFormationByIdAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	formation, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get formation"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(formation)
}

func UpdateFormationAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	existing, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find formation"})
		return
	}

	var f models.Formation
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	f.IdFormation = id
	if f.Status == "" {
		f.Status = existing.Status
	}

	if err := db.UpdateFormationAdmin(f); err != nil {
		fmt.Println("UpdateFormationAdmin error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update formation"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "formation mise à jour"})
}

func DeleteFormationAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	_, err := db.GetFormationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find formation"})
		return
	}

	if err := db.DeleteFormation(id); err != nil {
		fmt.Println("DeleteFormationAdmin error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete formation"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "formation supprimée"})
}
