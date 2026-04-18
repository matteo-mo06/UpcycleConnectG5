package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
)

func GetPublicAnnouncements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list, err := db.GetPublicAnnouncements(r.URL.Query().Get("search"), r.URL.Query().Get("id_category"))
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
	_ = json.NewEncoder(w).Encode(list)
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

	a := models.Announcement{
		Id_announcement:         uuid.New().String(),
		Id_category:             body.IdCategory,
		Title_announcement:      body.Title,
		Address_annoucement:     body.Address,
		City:                    body.City,
		Postal:                  body.Postal,
		Description_annoucement: body.Description,
		Availability_date:       body.AvailDate,
		Price:                   body.Price,
		Request:                 0,
		State_annoucement:       "Active",
		TypeAnnouncement:        typ,
		ConditionAnnouncement:   body.Condition,
	}

	if err := db.CreateUserAnnouncement(a, userID, body.PhotoURLs); err != nil {
		fmt.Println("CreateUserAnnouncement error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to create announcement"})
		return
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

	announcements, err := db.GetAllAnnouncements(idCategory, request)
	if err != nil {
		fmt.Println("GetAnnouncements error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch announcements"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(announcements)
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

	var a models.Announcement
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if a.Availability_date == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "availability_date is required"})
		return
	}

	a.Id_announcement = uuid.New().String()

	err = db.CreateAnnouncement(a)
	if err != nil {
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

	a.Id_announcement = id

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
