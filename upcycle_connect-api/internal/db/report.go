package db

import (
	"fmt"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func CreateReport(idReporter, idReportedUser, idAnnouncement, idTopic, idPost, reason string) error {
	var ann, topic, post, reported interface{}
	if idAnnouncement != "" {
		ann = idAnnouncement
	}
	if idTopic != "" {
		topic = idTopic
	}
	if idPost != "" {
		post = idPost
	}
	if idReportedUser != "" {
		reported = idReportedUser
	}
	_, err := config.Conn.Exec(
		`INSERT INTO REPORT (id_report, id_user, id_reported_user, id_announcement, id_topic, id_post, reason, status) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')`,
		uuid.New().String(), idReporter, reported, ann, topic, post, reason,
	)
	return err
}

func GetReportByID(id string) (models.Report, error) {
	var r models.Report
	err := config.Conn.QueryRow(`
		SELECT id_report, COALESCE(id_announcement, ''), COALESCE(id_topic, ''), COALESCE(id_post, ''), status
		FROM REPORT WHERE id_report = ?`, id,
	).Scan(&r.IdReport, &r.IdAnnouncement, &r.IdTopic, &r.IdPost, &r.Status)
	return r, err
}

func GetReports(status, search string) ([]models.Report, error) {
	query := `
		SELECT r.id_report, r.id_user, COALESCE(r.id_reported_user, ''), COALESCE(r.id_announcement, ''), COALESCE(r.id_topic, ''), COALESCE(r.id_post, ''),
		       r.reason, r.status, COALESCE(r.action_taken, ''), COALESCE(r.resolved_by, ''), COALESCE(r.resolved_at, ''), r.created_at,
		       TRIM(CONCAT(COALESCE(rep.first_name, ''), ' ', COALESCE(rep.last_name, ''))) AS reporter_name,
		       COALESCE(rep.email, '') AS reporter_email,
		       TRIM(CONCAT(COALESCE(repd.first_name, ''), ' ', COALESCE(repd.last_name, ''))) AS reported_name,
		       COALESCE(a.title_announcement, '') AS content_title
		FROM REPORT r
		LEFT JOIN USER rep  ON rep.id_user  = r.id_user
		LEFT JOIN USER repd ON repd.id_user = r.id_reported_user
		LEFT JOIN ANNOUNCEMENT a ON a.id_announcement = r.id_announcement
		WHERE 1=1`
	var args []any
	if status != "" {
		query += " AND r.status = ?"
		args = append(args, status)
	}
	if search != "" {
		query += " AND (TRIM(CONCAT(COALESCE(rep.first_name,''), ' ', COALESCE(rep.last_name,''))) LIKE ? OR COALESCE(a.title_announcement,'') LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	query += " ORDER BY r.created_at DESC"

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Report
	for rows.Next() {
		var r models.Report
		if err := rows.Scan(
			&r.IdReport, &r.IdReporter, &r.IdReportedUser, &r.IdAnnouncement, &r.IdTopic, &r.IdPost,
			&r.Reason, &r.Status, &r.ActionTaken, &r.ResolvedBy, &r.ResolvedAt, &r.CreatedAt,
			&r.ReporterName, &r.ReporterEmail, &r.ReportedUserName, &r.ContentTitle,
		); err != nil {
			return nil, err
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
		list = append(list, r)
	}
	return list, nil
}

func GetReportStats() (map[string]int, error) {
	stats := map[string]int{"total": 0, "pending": 0, "resolved": 0, "ignored": 0}
	var total int
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM REPORT").Scan(&total); err != nil {
		return nil, err
	}
	stats["total"] = total
	rows, err := config.Conn.Query("SELECT status, COUNT(*) FROM REPORT GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		var c int
		if err := rows.Scan(&s, &c); err != nil {
			return nil, err
		}
		stats[s] = c
	}
	return stats, nil
}

func UpdateReport(id, status, resolvedBy, actionTaken string) error {
	var rb, at interface{}
	if resolvedBy != "" {
		rb = resolvedBy
	}
	if actionTaken != "" {
		at = actionTaken
	}
	_, err := config.Conn.Exec(
		`UPDATE REPORT SET status = ?, resolved_by = ?, action_taken = ?, resolved_at = NOW() WHERE id_report = ?`,
		status, rb, at, id,
	)
	return err
}

func SoftDeleteContent(contentType, contentId string) error {
	var table, idCol string
	switch contentType {
	case "announcement":
		table = "ANNOUNCEMENT"
		idCol = "id_announcement"
	case "topic":
		table = "TOPIC"
		idCol = "id_topic"
	case "post":
		table = "POST"
		idCol = "id_post"
	default:
		return fmt.Errorf("type de contenu inconnu: %s", contentType)
	}
	_, err := config.Conn.Exec("UPDATE "+table+" SET deleted_at = NOW() WHERE "+idCol+" = ?", contentId)
	return err
}
