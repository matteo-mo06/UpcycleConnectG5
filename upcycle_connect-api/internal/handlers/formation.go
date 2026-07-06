package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
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

func GetFormationDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	docs, err := db.GetDocumentsByCategory(id)
	if err != nil {
		fmt.Println("GetFormationDocuments error:", err)
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

func createFormation(w http.ResponseWriter, r *http.Request, defaultStatus string) {
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

	if req.Address == nil || *req.Address == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "l'adresse est requise"})
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
		AddressFormation:     req.Address,
		CityFormation:        req.City,
		PostalFormation:      req.Postal,
		Capacity:             req.Capacity,
		Level:                req.Level,
		DurationHours:        req.DurationH,
		IdCategory:           req.IdCategory,
		Price:                req.Price,
		Syllabus:             req.Syllabus,
		Prerequisites:        req.Prerequisites,
		Objectives:           req.Objectives,
		Status:               defaultStatus,
		IdCreator:            &userID,
		IdFormateur:          &userID,
	}

	if err := db.CreateFormation(f); err != nil {
		fmt.Println("CreateFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create formation"})
		return
	}

	for _, url := range req.PhotoURLs {
		if err := db.CreateDocument(models.Document{
			IdDocument: uuid.New().String(),
			IdUser:     userID,
			Category:   f.IdFormation,
			Link:       url,
		}); err != nil {
			fmt.Println("CreateFormation document error:", err)
		}
	}

	msg := "formation créée, en attente de validation"
	if defaultStatus == "approved" {
		msg = "formation créée"
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": msg, "formation": f})
}

func CreateFormation(w http.ResponseWriter, r *http.Request) {
	createFormation(w, r, "pending")
}

func CreateFormationAdmin(w http.ResponseWriter, r *http.Request) {
	createFormation(w, r, "approved")
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

	if req.Address == nil || *req.Address == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "l'adresse est requise"})
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
	formation.AddressFormation = req.Address
	formation.CityFormation = req.City
	formation.PostalFormation = req.Postal
	formation.Capacity = req.Capacity
	formation.Level = req.Level
	formation.DurationHours = req.DurationH
	formation.IdCategory = req.IdCategory
	formation.Price = req.Price
	formation.Syllabus = req.Syllabus
	formation.Prerequisites = req.Prerequisites
	formation.Objectives = req.Objectives

	if err := db.UpdateFormation(formation); err != nil {
		fmt.Println("UpdateMyFormation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update formation"})
		return
	}

	for _, url := range req.PhotoURLs {
		if err := db.CreateDocument(models.Document{
			IdDocument: uuid.New().String(),
			IdUser:     userID,
			Category:   id,
			Link:       url,
		}); err != nil {
			fmt.Println("UpdateMyFormation document error:", err)
		}
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

func GetDeletedFormations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	formations, err := db.GetDeletedFormations()
	if err != nil {
		fmt.Println("GetDeletedFormations error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch formations"})
		return
	}

	if formations == nil {
		formations = []models.DeletedFormation{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(formations)
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

func GetFormationSyllabusPDF(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	perms, _ := r.Context().Value(middleware.ContextPermissions).([]string)

	formation, err := db.GetFormationById(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formation not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get formation"})
		return
	}

	canSee := formation.Status == "approved" ||
		(formation.IdCreator != nil && *formation.IdCreator == userID) ||
		slices.Contains(perms, "manage_formations")

	if !canSee {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès non autorisé"})
		return
	}

	pdfBytes, err := buildSyllabusPDF(formation)
	if err != nil {
		fmt.Println("GetFormationSyllabusPDF error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur génération PDF"})
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="syllabus-%s.pdf"`, formation.IdFormation))
	_, _ = w.Write(pdfBytes)
}

func buildSyllabusPDF(f models.Formation) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 15, 20)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetTextColor(45, 106, 79)
	pdf.SetFont("Helvetica", "B", 20)
	pdf.Cell(0, 10, "UpcycleConnect")
	pdf.Ln(9)
	pdf.SetTextColor(120, 120, 120)
	pdf.SetFont("Helvetica", "", 9)
	pdf.Cell(0, 5, tr("Programme de formation"))
	pdf.Ln(10)

	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(7)

	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Helvetica", "B", 16)
	pdf.MultiCell(170, 8, tr(f.TitleFormation), "", "L", false)

	formateur := ""
	if f.FormateurName != nil && *f.FormateurName != "" {
		formateur = *f.FormateurName
	} else if f.CreatorName != nil {
		formateur = *f.CreatorName
	}
	if formateur != "" {
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(90, 90, 90)
		pdf.Cell(0, 6, tr("Animée par "+formateur))
		pdf.Ln(9)
	} else {
		pdf.Ln(3)
	}

	levelLabels := map[string]string{"beginner": "Débutant", "intermediate": "Intermédiaire", "advanced": "Avancé"}
	levelLabel := levelLabels[f.Level]
	if levelLabel == "" {
		levelLabel = f.Level
	}

	dateStr := "-"
	if f.DateFormation != nil && *f.DateFormation != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", *f.DateFormation); err == nil {
			dateStr = parsed.Format("02/01/2006 à 15:04")
		} else if parsed, err := time.Parse("2006-01-02T15:04:05", *f.DateFormation); err == nil {
			dateStr = parsed.Format("02/01/2006 à 15:04")
		} else {
			dateStr = *f.DateFormation
		}
	}

	location := "-"
	var parts []string
	if f.AddressFormation != nil && *f.AddressFormation != "" {
		parts = append(parts, *f.AddressFormation)
	}
	cityLine := ""
	if f.PostalFormation != nil && *f.PostalFormation != "" {
		cityLine = *f.PostalFormation
	}
	if f.CityFormation != nil && *f.CityFormation != "" {
		if cityLine != "" {
			cityLine += " "
		}
		cityLine += *f.CityFormation
	}
	if cityLine != "" {
		parts = append(parts, cityLine)
	}
	if len(parts) > 0 {
		location = strings.Join(parts, ", ")
	}

	duration := "-"
	if f.DurationHours != nil {
		duration = fmt.Sprintf("%dh", *f.DurationHours)
	}

	price := "Gratuit"
	if f.Price != nil && *f.Price > 0 {
		price = fmt.Sprintf("%.2f EUR", *f.Price)
	}

	infoRow := func(label, value string) {
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetTextColor(40, 40, 40)
		pdf.CellFormat(30, 7, tr(label), "", 0, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 9)
		pdf.SetTextColor(70, 70, 70)
		pdf.CellFormat(0, 7, tr(value), "", 1, "L", false, 0, "")
	}
	infoRow("Niveau :", levelLabel)
	infoRow("Date :", dateStr)
	infoRow("Lieu :", location)
	infoRow("Durée :", duration)
	infoRow("Prix :", price)

	pdf.Ln(3)

	section := func(heading string, content *string) {
		if content == nil || *content == "" {
			return
		}
		pdf.SetFont("Helvetica", "B", 11)
		pdf.SetTextColor(45, 106, 79)
		pdf.Cell(0, 8, tr(heading))
		pdf.Ln(7)
		pdf.SetFont("Helvetica", "", 9.5)
		pdf.SetTextColor(60, 60, 60)
		pdf.MultiCell(170, 5.5, tr(*content), "", "L", false)
		pdf.Ln(4)
	}

	section("Description", f.DescriptionFormation)
	section("Objectifs pédagogiques", f.Objectives)
	section("Prérequis", f.Prerequisites)
	section("Programme détaillé", f.Syllabus)

	if steps, err := db.GetFormationSteps(f.IdFormation); err == nil && len(steps) > 0 {
		pdf.SetFont("Helvetica", "B", 11)
		pdf.SetTextColor(45, 106, 79)
		pdf.Cell(0, 8, tr("Programme / Étapes"))
		pdf.Ln(7)
		for i, s := range steps {
			pdf.SetFont("Helvetica", "B", 9.5)
			pdf.SetTextColor(60, 60, 60)
			pdf.MultiCell(170, 5.5, tr(fmt.Sprintf("%d. %s", i+1, s.Title)), "", "L", false)
			if s.Description != nil && *s.Description != "" {
				pdf.SetFont("Helvetica", "", 9.5)
				pdf.SetTextColor(60, 60, 60)
				pdf.MultiCell(170, 5.5, tr(*s.Description), "", "L", false)
			}
		}
		pdf.Ln(4)
	}

	pdf.SetY(-25)
	pdf.SetTextColor(120, 120, 120)
	pdf.SetFont("Helvetica", "", 7.5)
	pdf.CellFormat(0, 4, tr("UpcycleConnect - contact@upcycleconnect.fr - upcycleconnect.fr"), "", 0, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
