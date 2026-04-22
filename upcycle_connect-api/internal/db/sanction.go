package db

import (
	"time"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func CreateSanction(idUser, idAdmin, idReport, sanctionType, reason string, durationDays int) error {
	var idReportVal interface{}
	if idReport != "" {
		idReportVal = idReport
	}
	var expiresAt interface{}
	if durationDays > 0 {
		expiresAt = time.Now().AddDate(0, 0, durationDays).Format("2006-01-02 15:04:05")
	}
	_, err := config.Conn.Exec(
		`INSERT INTO USER_SANCTIONS (id_sanction, id_user, id_admin, id_report, type, reason, duration_days, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		uuid.New().String(), idUser, idAdmin, idReportVal, sanctionType, reason, durationDays, expiresAt,
	)
	return err
}

func GetUserSanctions(userId string) ([]models.Sanction, error) {
	rows, err := config.Conn.Query(`
		SELECT s.id_sanction, s.id_user, s.id_admin, COALESCE(s.id_report, ''), s.type,
		       COALESCE(s.reason, ''), COALESCE(s.duration_days, 0), COALESCE(s.expires_at, ''), s.created_at,
		       TRIM(CONCAT(COALESCE(a.first_name, ''), ' ', COALESCE(a.last_name, ''))) AS admin_name
		FROM USER_SANCTIONS s
		LEFT JOIN USER a ON a.id_user = s.id_admin
		WHERE s.id_user = ?
		ORDER BY s.created_at DESC`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Sanction
	for rows.Next() {
		var s models.Sanction
		if err := rows.Scan(
			&s.IdSanction, &s.IdUser, &s.IdAdmin, &s.IdReport, &s.Type,
			&s.Reason, &s.DurationDays, &s.ExpiresAt, &s.CreatedAt, &s.AdminName,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func GetUserHistorySummary(search string) ([]models.UserHistorySummary, error) {
	query := `
		SELECT u.id_user,
		       TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS name,
		       u.email,
		       COUNT(DISTINCT r.id_report) AS report_count,
		       COUNT(DISTINCT s.id_sanction) AS sanction_count
		FROM USER u
		LEFT JOIN REPORT r ON r.id_reported_user = u.id_user
		LEFT JOIN USER_SANCTIONS s ON s.id_user = u.id_user
		WHERE 1=1`
	var args []any
	if search != "" {
		query += " AND (u.first_name LIKE ? OR u.last_name LIKE ? OR u.email LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	query += " GROUP BY u.id_user, u.first_name, u.last_name, u.email"
	query += " HAVING report_count > 0 OR sanction_count > 0"
	query += " ORDER BY report_count DESC, sanction_count DESC"

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.UserHistorySummary
	for rows.Next() {
		var u models.UserHistorySummary
		if err := rows.Scan(&u.IdUser, &u.Name, &u.Email, &u.ReportCount, &u.SanctionCount); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

func GetUserHistory(userId string) (models.UserHistory, error) {
	var h models.UserHistory
	sanctions, err := GetUserSanctions(userId)
	if err != nil {
		return h, err
	}
	h.Sanctions = sanctions

	rows, err := config.Conn.Query(`
		SELECT r.id_report, COALESCE(r.id_announcement, ''), COALESCE(r.id_topic, ''), COALESCE(r.id_post, ''),
		       r.reason, r.status, r.created_at,
		       COALESCE(a.title_announcement, '') AS content_title
		FROM REPORT r
		LEFT JOIN ANNOUNCEMENT a ON a.id_announcement = r.id_announcement
		WHERE r.id_reported_user = ?
		ORDER BY r.created_at DESC`, userId)
	if err != nil {
		return h, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Report
		if err := rows.Scan(
			&r.IdReport, &r.IdAnnouncement, &r.IdTopic, &r.IdPost,
			&r.Reason, &r.Status, &r.CreatedAt, &r.ContentTitle,
		); err != nil {
			return h, err
		}
		switch {
		case r.IdAnnouncement != "":
			r.ContentType = "announcement"
		case r.IdTopic != "":
			r.ContentType = "topic"
			if r.ContentTitle == "" {
				r.ContentTitle = r.IdTopic
			}
		case r.IdPost != "":
			r.ContentType = "post"
			if r.ContentTitle == "" {
				r.ContentTitle = r.IdPost
			}
		}
		h.ReportsReceived = append(h.ReportsReceived, r)
	}
	return h, nil
}
