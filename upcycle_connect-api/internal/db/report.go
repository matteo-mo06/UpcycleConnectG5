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

func GetReports(status, search string, limit, offset int) ([]models.Report, int, error) {
	where := `WHERE 1=1`
	var args []any
	if status != "" {
		where += " AND r.status = ?"
		args = append(args, status)
	}
	if search != "" {
		where += " AND (TRIM(CONCAT(COALESCE(rep.first_name,''), ' ', COALESCE(rep.last_name,''))) LIKE ? OR COALESCE(a.title_announcement,'') LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := config.Conn.QueryRow(`SELECT COUNT(*) FROM REPORT r
		LEFT JOIN USER rep ON rep.id_user = r.id_user
		LEFT JOIN ANNOUNCEMENT a ON a.id_announcement = r.id_announcement
		`+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT r.id_report, r.id_user, COALESCE(r.id_reported_user, ''), COALESCE(r.id_announcement, ''), COALESCE(r.id_topic, ''), COALESCE(r.id_post, ''),
		       r.reason, r.status, COALESCE(r.action_taken, ''), COALESCE(r.resolved_by, ''), COALESCE(r.resolved_at, ''), r.created_at,
		       TRIM(CONCAT(COALESCE(rep.first_name, ''), ' ', COALESCE(rep.last_name, ''))) AS reporter_name,
		       COALESCE(rep.email, '') AS reporter_email,
		       TRIM(CONCAT(COALESCE(repd.first_name, ''), ' ', COALESCE(repd.last_name, ''))) AS reported_name,
		       COALESCE(a.title_announcement, t.title_topic, LEFT(rp.body_post, 80), '') AS content_title,
		       COALESCE(LEFT(rp.body_post, 300), (SELECT LEFT(pp.body_post, 300) FROM POST pp JOIN TOPIC_POST tp ON tp.id_post = pp.id_post WHERE tp.id_topic = r.id_topic AND pp.deleted_at IS NULL ORDER BY pp.created_at ASC LIMIT 1), '') AS content_body,
		       COALESCE((SELECT d.link FROM DOCUMENT d WHERE d.category = r.id_announcement ORDER BY d.id_document LIMIT 1), '') AS content_photo
		FROM REPORT r
		LEFT JOIN USER rep  ON rep.id_user  = r.id_user
		LEFT JOIN USER repd ON repd.id_user = r.id_reported_user
		LEFT JOIN ANNOUNCEMENT a ON a.id_announcement = r.id_announcement
		LEFT JOIN TOPIC t ON t.id_topic = r.id_topic
		LEFT JOIN POST rp ON rp.id_post = r.id_post
		`+where+` ORDER BY r.created_at DESC LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Report
	for rows.Next() {
		var r models.Report
		if err := rows.Scan(
			&r.IdReport, &r.IdReporter, &r.IdReportedUser, &r.IdAnnouncement, &r.IdTopic, &r.IdPost,
			&r.Reason, &r.Status, &r.ActionTaken, &r.ResolvedBy, &r.ResolvedAt, &r.CreatedAt,
			&r.ReporterName, &r.ReporterEmail, &r.ReportedUserName, &r.ContentTitle, &r.ContentBody, &r.ContentPhoto,
		); err != nil {
			return nil, 0, err
		}
		switch {
		case r.IdAnnouncement != "":
			r.ContentType = "announcement"
		case r.IdTopic != "":
			r.ContentType = "topic"
		case r.IdPost != "":
			r.ContentType = "post"
		}
		list = append(list, r)
	}
	return list, total, nil
}

func GetReportStats() (map[string]int, error) {
	stats := map[string]int{"total": 0, "pending": 0, "resolved": 0}
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
