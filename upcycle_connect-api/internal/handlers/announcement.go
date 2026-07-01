package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/google/uuid"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
	"upcycle_connect-api/internal/utils"
)

func GetPublicAnnouncements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, limit, offset := parsePage(r, 12)
	list, total, err := db.GetPublicAnnouncements(r.URL.Query().Get("search"), r.URL.Query().Get("id_category"), r.URL.Query().Get("type"), limit, offset)
	if err != nil {
		fmt.Println("GetPublicAnnouncements error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch announcements"})
		return
	}
	if list == nil {
		list = []models.Announcement{}
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(list, total, page, limit))
}

func GetPublicAnnouncementById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing announcement id"})
		return
	}

	announcement, err := db.GetAnnouncementById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "announcement not found"})
			return
		}
		fmt.Println("GetPublicAnnouncementById error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get announcement"})
		return
	}

	docs, _ := db.GetDocumentsByCategory(id)
	if docs == nil {
		docs = []models.Document{}
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"announcement": announcement, "photos": docs})
}

func CreateUserAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var body models.CreateAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
	if body.Title == "" || body.AvailDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "title et availability_date requis"})
		return
	}
	if body.Type == "don" {
		body.Price = 0
	}
	if body.Type == "vente" && body.Price < 0.01 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "le prix minimum pour une vente est 0.01 €"})
		return
	}

	typ := body.Type
	if typ == "" {
		typ = "don"
	}

	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)

	if slices.Contains(roles, config.RoleArtisan) && !db.IsUserPremium(userID) {
		count, _ := db.GetUserAnnouncementCount(userID)
		if count >= 3 {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "limite de 3 annonces atteinte, passez à l'abonnement premium pour en publier davantage"})
			return
		}
	}

	state := "En attente"
	if slices.Contains(roles, config.RoleAdmin) {
		state = "Active"
	}

	a := models.Announcement{
		IdAnnouncement:          uuid.New().String(),
		IdCategory:              body.IdCategory,
		TitleAnnouncement:       body.Title,
		AddressAnnouncement:     body.Address,
		City:                    body.City,
		Postal:                  body.Postal,
		DescriptionAnnouncement: body.Description,
		AvailabilityDate:        body.AvailDate,
		Price:                   body.Price,
		Request:                 0,
		StateAnnouncement:       state,
		TypeAnnouncement:        typ,
		ConditionAnnouncement:   body.Condition,
	}

	if err := db.CreateUserAnnouncement(a, userID, body.PhotoURLs); err != nil {
		fmt.Println("CreateUserAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create announcement"})
		return
	}

	if slices.Contains(roles, config.RoleArtisan) {
		_ = db.IncrementUserAnnouncementCount(userID)
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "annonce créée", "announcement": a})
}

func GetMyAnnouncements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	list, err := db.GetUserAnnouncements(userID)
	if err != nil {
		fmt.Println("GetMyAnnouncements error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch announcements"})
		return
	}
	if list == nil {
		list = []models.Announcement{}
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func DeleteMyAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	owner, err := db.IsAnnouncementOwner(id, userID)
	if err != nil {
		fmt.Println("DeleteMyAnnouncement ownership check error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	if !owner {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous n'êtes pas propriétaire de cette annonce"})
		return
	}

	if err := db.DeleteAnnouncement(id); err != nil {
		fmt.Println("DeleteMyAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete announcement"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "annonce supprimée"})
}

func DeleteAnnouncementWithPermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	if err := db.DeleteAnnouncement(id); err != nil {
		fmt.Println("DeleteAnnouncementWithPermission error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete announcement"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "annonce supprimée"})
}

func GetAnnouncementDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing announcement id"})
		return
	}

	docs, err := db.GetDocumentsByCategory(id)
	if err != nil {
		fmt.Println("GetAnnouncementDocuments error:", err)
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

func ApproveAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	if err := db.SetAnnouncementState(id, "Active"); err != nil {
		fmt.Println("ApproveAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to approve announcement"})
		return
	}
	go func() {
		if ownerID, err := db.GetAnnouncementOwnerID(id); err == nil {
			utils.SendPushNotification(db.GetOnesignalPlayerID(ownerID), "Annonce approuvée", "Votre annonce a été approuvée et est maintenant visible.")
		}
	}()
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "annonce approuvée"})
}

func RejectAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	var body struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if err := db.RejectAnnouncementWithReason(id, body.Reason); err != nil {
		fmt.Println("RejectAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to reject announcement"})
		return
	}
	go func() {
		if ownerID, err := db.GetAnnouncementOwnerID(id); err == nil {
			utils.SendPushNotification(db.GetOnesignalPlayerID(ownerID), "Annonce refusée", "Votre annonce n'a pas été approuvée.")
		}
	}()
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "annonce rejetée"})
}

func GetAnnouncementStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := db.GetAnnouncementStats()
	if err != nil {
		fmt.Println("GetAnnouncementStats error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch announcement stats"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}

func GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idCategory := r.URL.Query().Get("id_category")
	search := r.URL.Query().Get("search")
	filterType := r.URL.Query().Get("type")
	filterStatus := r.URL.Query().Get("status")

	request := -1
	if raw := r.URL.Query().Get("request"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "request must be 0 or 1"})
			return
		}
		request = v
	}

	page, limit, offset := parsePage(r, 20)

	announcements, total, err := db.GetAllAnnouncements(idCategory, search, filterType, filterStatus, request, limit, offset)
	if err != nil {
		fmt.Println("GetAnnouncements error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch announcements"})
		return
	}
	if announcements == nil {
		announcements = []models.Announcement{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pageResponse(announcements, total, page, limit))
}

func GetAnnouncementById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing announcement id"})
		return
	}

	announcement, err := db.GetAnnouncementById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "announcement not found"})
			return
		}
		fmt.Println("GetAnnouncementById error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get announcement"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(announcement)
}

func CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var a models.Announcement
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if a.AvailabilityDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "availability_date is required"})
		return
	}

	a.IdAnnouncement = uuid.New().String()

	if err = db.CreateUserAnnouncement(a, userID, nil); err != nil {
		fmt.Println("CreateAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create announcement"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "announcement created successfully", "announcement": a})
}

func UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing announcement id"})
		return
	}

	_, err := db.GetAnnouncementById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "announcement not found"})
			return
		}
		fmt.Println("UpdateAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find announcement"})
		return
	}

	var a models.Announcement
	err = json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	a.IdAnnouncement = id

	err = db.UpdateAnnouncement(a)
	if err != nil {
		fmt.Println("UpdateAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update announcement"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "announcement updated successfully", "announcement": a})
}

func UpdateMyAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	ann, err := db.GetAnnouncementById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "annonce introuvable"})
			return
		}
		fmt.Println("UpdateMyAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find announcement"})
		return
	}

	ownerID, err := db.GetAnnouncementOwnerID(id)
	if err != nil || ownerID != userID {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous n'êtes pas propriétaire de cette annonce"})
		return
	}

	if ann.StateAnnouncement == "Vendu" || ann.StateAnnouncement == "Supprimée" {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette annonce ne peut plus être modifiée"})
		return
	}

	var body models.Announcement
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	ann.TitleAnnouncement = body.TitleAnnouncement
	ann.DescriptionAnnouncement = body.DescriptionAnnouncement
	ann.Price = body.Price
	ann.AvailabilityDate = body.AvailabilityDate
	ann.ConditionAnnouncement = body.ConditionAnnouncement
	ann.TypeAnnouncement = body.TypeAnnouncement
	ann.AddressAnnouncement = body.AddressAnnouncement
	ann.City = body.City
	ann.Postal = body.Postal
	if body.IdCategory != "" {
		ann.IdCategory = body.IdCategory
	}

	if err := db.UpdateAnnouncement(ann); err != nil {
		fmt.Println("UpdateMyAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to update announcement"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "annonce mise à jour"})
}

func ClaimAnnouncement(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	roles, _ := r.Context().Value(middleware.ContextRoles).([]string)
	if slices.Contains(roles, config.RoleArtisan) && !db.IsUserPremium(userID) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "abonnement premium requis pour acquérir des annonces"})
		return
	}

	owner, err := db.IsAnnouncementOwner(id, userID)
	if err != nil {
		http.Error(w, "erreur serveur", http.StatusInternalServerError)
		return
	}
	if owner {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous ne pouvez pas acquérir votre propre annonce"})
		return
	}

	ann, err := db.GetAnnouncementById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "annonce introuvable"})
		return
	}
	if ann.Price > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette annonce est payante, utilisez le système de paiement"})
		return
	}

	if err := db.ClaimAnnouncement(id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "cette annonce a déjà été acquise"})
			return
		}
		fmt.Println("ClaimAnnouncement error:", err)
		http.Error(w, "erreur serveur", http.StatusInternalServerError)
		return
	}

	if err := db.AwardScore(userID, "announcement_bought", id); err != nil {
		fmt.Println("AwardScore announcement_bought error:", err)
	}
	if sellerID, err := db.GetAnnouncementOwnerID(id); err == nil {
		if err := db.AwardScore(sellerID, "announcement_sold", id); err != nil {
			fmt.Println("AwardScore announcement_sold error:", err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetMyAcquisitions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)

	list, err := db.GetUserAcquisitions(userID)
	if err != nil {
		fmt.Println("GetMyAcquisitions error:", err)
		http.Error(w, "erreur serveur", http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = []models.Announcement{}
	}
	_ = json.NewEncoder(w).Encode(list)
}

func DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing announcement id"})
		return
	}

	_, err := db.GetAnnouncementById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "announcement not found"})
			return
		}
		fmt.Println("DeleteAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to find announcement"})
		return
	}

	err = db.DeleteAnnouncement(id)
	if err != nil {
		fmt.Println("DeleteAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete announcement"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "announcement deleted successfully"})
}

func GetMyInvoice(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	announcementID := r.PathValue("id")

	isBuyer, err := db.IsAnnouncementBuyer(announcementID, userID)
	if err != nil || !isBuyer {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}

	invoicePath, err := db.GetInvoicePath(announcementID)
	if err != nil || invoicePath == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "facture non disponible"})
		return
	}

	fullPath := filepath.Join(config.InvoicesDir(), invoicePath)
	f, err := os.Open(fullPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "fichier introuvable"})
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="facture.pdf"`)
	_, _ = io.Copy(w, f)
}

func RequestFeature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	if !db.IsUserPremium(userID) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "abonnement premium requis"})
		return
	}

	status, slot, err := db.RequestFeature(id, userID)
	if err != nil {
		switch err {
		case db.ErrNotOwner:
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		case db.ErrNotActive:
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "l'annonce doit être active"})
		case db.ErrMonthlyQuota:
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "quota mensuel atteint"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		}
		return
	}

	resp := map[string]any{"status": status}
	if slot != nil {
		if status == "active" {
			resp["featured_until"] = slot.Format("2006-01-02T15:04:05Z")
		} else {
			resp["estimated_slot"] = slot.Format("2006-01-02T15:04:05Z")
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func ConfirmFeature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	status, slot, err := db.ConfirmFeatureRequest(id, userID)
	if err != nil {
		switch err {
		case db.ErrNotOwner:
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		case db.ErrNotQueued:
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "aucune demande en file d'attente"})
		case db.ErrNoSlot:
			resp := map[string]any{"error": "pas encore de slot disponible"}
			if slot != nil {
				resp["estimated_slot"] = slot.Format("2006-01-02T15:04:05Z")
			}
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(resp)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		}
		return
	}

	resp := map[string]any{"status": status}
	if slot != nil {
		resp["featured_until"] = slot.Format("2006-01-02T15:04:05Z")
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func CancelFeature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	ok, err := db.IsAnnouncementOwner(id, userID)
	if err != nil || !ok {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès refusé"})
		return
	}

	if err := db.CancelFeatureRequest(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "demande annulée"})
}

func AdminToggleFeature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	var body struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps invalide"})
		return
	}

	if err := db.AdminToggleFeature(id, body.Active); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "mise en avant mise à jour"})
}