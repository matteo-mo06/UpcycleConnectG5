package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"upcycle_connect-api/internal/db"
	"upcycle_connect-api/internal/middleware"
	"upcycle_connect-api/internal/models"
)

func SubmitReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reporterID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var body struct {
		IdAnnouncement string `json:"id_announcement"`
		IdTopic        string `json:"id_topic"`
		IdPost         string `json:"id_post"`
		Reason         string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}

	count := 0
	if body.IdAnnouncement != "" {
		count++
	}
	if body.IdTopic != "" {
		count++
	}
	if body.IdPost != "" {
		count++
	}
	if count != 1 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "un seul contenu peut être signalé à la fois"})
		return
	}
	if body.Reason == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "la raison est requise"})
		return
	}

	var reportedUserID string
	type contentLookup struct {
		id     string
		fetch  func(string) (string, error)
		errMsg string
	}
	for _, l := range []contentLookup{
		{body.IdAnnouncement, db.GetAnnouncementOwnerID, "annonce introuvable"},
		{body.IdTopic, db.GetTopicAuthor, "topic introuvable"},
		{body.IdPost, db.GetPostAuthor, "post introuvable"},
	} {
		if l.id == "" {
			continue
		}
		authorID, err := l.fetch(l.id)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": l.errMsg})
				return
			}
			fmt.Println("SubmitReport author lookup error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
			return
		}
		if authorID == reporterID {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "vous ne pouvez pas signaler votre propre contenu"})
			return
		}
		reportedUserID = authorID
		break
	}

	if err := db.CreateReport(reporterID, reportedUserID, body.IdAnnouncement, body.IdTopic, body.IdPost, body.Reason); err != nil {
		fmt.Println("SubmitReport error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création du signalement"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "signalement envoyé"})
}

func GetAdminReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, limit, offset := parsePage(r, 20)
	list, total, err := db.GetReports(r.URL.Query().Get("status"), r.URL.Query().Get("search"), limit, offset)
	if err != nil {
		fmt.Println("GetAdminReports error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la récupération des signalements"})
		return
	}
	if list == nil {
		list = []models.Report{}
	}
	_ = json.NewEncoder(w).Encode(pageResponse(list, total, page, limit))
}

func GetAdminReportStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := db.GetReportStats()
	if err != nil {
		fmt.Println("GetAdminReportStats error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la récupération des stats"})
		return
	}
	_ = json.NewEncoder(w).Encode(stats)
}

func UpdateAdminReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	adminID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := r.PathValue("id")

	var body struct {
		Status      string `json:"status"`
		ActionTaken string `json:"action_taken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}

	allowed := map[string]bool{"resolved": true}
	if !allowed[body.Status] {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "status invalide (resolved)"})
		return
	}

	if err := db.UpdateReport(id, body.Status, adminID, body.ActionTaken); err != nil {
		fmt.Println("UpdateAdminReport error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la mise à jour"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "signalement mis à jour"})
}

func SoftDeleteReportContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	adminID, _ := r.Context().Value(middleware.ContextUserID).(string)
	reportID := r.PathValue("id")

	report, err := db.GetReportByID(reportID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "signalement introuvable"})
			return
		}
		fmt.Println("SoftDeleteReportContent GetReportByID error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur serveur"})
		return
	}

	var contentType, contentId string
	switch {
	case report.IdAnnouncement != "":
		contentType = "announcement"
		contentId = report.IdAnnouncement
	case report.IdTopic != "":
		contentType = "topic"
		contentId = report.IdTopic
	case report.IdPost != "":
		contentType = "post"
		contentId = report.IdPost
	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "aucun contenu associé à ce signalement"})
		return
	}

	if err := db.SoftDeleteContent(contentType, contentId); err != nil {
		fmt.Println("SoftDeleteReportContent error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la suppression du contenu"})
		return
	}
	_ = db.UpdateReport(reportID, "resolved", adminID, "content_removed")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "contenu supprimé"})
}

func CreateSanction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	adminID, _ := r.Context().Value(middleware.ContextUserID).(string)
	userID := r.PathValue("id")

	var body struct {
		Type         string `json:"type"`
		Reason       string `json:"reason"`
		DurationDays int    `json:"duration_days"`
		IdReport     string `json:"id_report"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "corps de requête invalide"})
		return
	}

	allowed := map[string]bool{"warning": true, "suspension": true, "ban": true}
	if !allowed[body.Type] {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "type invalide (warning | suspension | ban)"})
		return
	}

	if err := db.CreateSanction(userID, adminID, body.IdReport, body.Type, body.Reason, body.DurationDays); err != nil {
		fmt.Println("CreateSanction error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la création de la sanction"})
		return
	}

	defaultMessages := map[string]string{
		"warning":    "Vous avez reçu un avertissement de la part des modérateurs.",
		"suspension": "Votre compte a été suspendu.",
		"ban":        "Votre compte a été banni.",
	}
	msg := body.Reason
	if msg == "" {
		msg = defaultMessages[body.Type]
	}
	if body.Type == "ban" {
		if err := db.UpdateUserStatus(userID, "blacklisted"); err != nil {
			fmt.Println("CreateSanction ban UpdateUserStatus error:", err)
		}
	} else if body.Type == "suspension" {
		if err := db.UpdateUserStatus(userID, "suspended"); err != nil {
			fmt.Println("CreateSanction suspend UpdateUserStatus error:", err)
		}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "sanction appliquée"})
}

func GetUserHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.PathValue("id")

	history, err := db.GetUserHistory(userID)
	if err != nil {
		fmt.Println("GetUserHistory error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la récupération de l'historique"})
		return
	}
	if history.Sanctions == nil {
		history.Sanctions = []models.Sanction{}
	}
	if history.ReportsReceived == nil {
		history.ReportsReceived = []models.Report{}
	}
	_ = json.NewEncoder(w).Encode(history)
}

func GetUsersHistorySummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list, err := db.GetUserHistorySummary(r.URL.Query().Get("search"))
	if err != nil {
		fmt.Println("GetUsersHistorySummary error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "erreur lors de la récupération"})
		return
	}
	if list == nil {
		list = []models.UserHistorySummary{}
	}
	_ = json.NewEncoder(w).Encode(list)
}
