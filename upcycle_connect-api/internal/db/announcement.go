package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

const announcementSelect = `SELECT a.id_announcement, a.id_category, a.title_announcement, a.address_annoucement, a.city, a.postal,
	a.description_annoucement, a.availability_date, a.price, a.request, a.state_annoucement,
	u.id_user, u.first_name, u.last_name, a.type_announcement, a.condition_announcement,
	(SELECT d.link FROM DOCUMENT d WHERE d.category = a.id_announcement ORDER BY d.id_document LIMIT 1) AS first_photo `

func scanAnnouncement(row interface{ Scan(...any) error }) (models.Announcement, error) {
	var a models.Announcement
	var idCat, authorId, firstName, lastName, typ, cond, firstPhoto sql.NullString
	err := row.Scan(
		&a.Id_announcement, &idCat, &a.Title_announcement,
		&a.Address_annoucement, &a.City, &a.Postal,
		&a.Description_annoucement, &a.Availability_date, &a.Price,
		&a.Request, &a.State_annoucement,
		&authorId, &firstName, &lastName, &typ, &cond, &firstPhoto,
	)
	if err != nil {
		return a, err
	}
	a.Id_category = idCat.String
	a.AuthorId = authorId.String
	a.AuthorName = strings.TrimSpace(firstName.String + " " + lastName.String)
	a.TypeAnnouncement = typ.String
	a.ConditionAnnouncement = cond.String
	a.FirstPhoto = firstPhoto.String
	return a, nil
}

func GetAllAnnouncements(idCategory string, request int) ([]models.Announcement, error) {
	query := announcementSelect + `FROM ANNOUNCEMENT a
	          LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
	          LEFT JOIN USER u ON u.id_user = ua.id_user
	          WHERE 1=1`
	var args []any

	if idCategory != "" {
		query += " AND a.id_category = ?"
		args = append(args, idCategory)
	}
	if request >= 0 {
		query += " AND a.request = ?"
		args = append(args, request)
	}

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Announcement
	for rows.Next() {
		a, err := scanAnnouncement(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func GetPublicAnnouncements(search, idCategory string) ([]models.Announcement, error) {
	query := announcementSelect + `FROM ANNOUNCEMENT a
	          LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
	          LEFT JOIN USER u ON u.id_user = ua.id_user
	          WHERE a.state_annoucement = 'Active' AND a.request = 0`
	var args []any

	if search != "" {
		query += " AND (a.title_announcement LIKE ? OR a.description_annoucement LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	if idCategory != "" {
		query += " AND a.id_category = ?"
		args = append(args, idCategory)
	}
	query += " ORDER BY a.availability_date DESC"

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Announcement
	for rows.Next() {
		a, err := scanAnnouncement(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func GetUserAnnouncements(userID string) ([]models.Announcement, error) {
	rows, err := config.Conn.Query(announcementSelect+`FROM ANNOUNCEMENT a
		JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		JOIN USER u ON u.id_user = ua.id_user
		WHERE ua.id_user = ?
		ORDER BY a.availability_date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Announcement
	for rows.Next() {
		a, err := scanAnnouncement(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}

	for i := range list {
		var num sql.NullInt64
		_ = config.Conn.QueryRow(`
			SELECT l.locker_number FROM ANNOUNCEMENT_LOCKER al
			JOIN LOCKER l ON l.id_locker = al.id_locker
			WHERE al.id_announcement = ?`, list[i].Id_announcement,
		).Scan(&num)
		list[i].LockerNumber = int(num.Int64)
	}

	return list, nil
}

func IsAnnouncementOwner(announcementID, userID string) (bool, error) {
	var count int
	err := config.Conn.QueryRow(
		"SELECT COUNT(*) FROM USER_ANNOUNCEMENT WHERE id_announcement = ? AND id_user = ?",
		announcementID, userID,
	).Scan(&count)
	return count > 0, err
}

func GetAnnouncementStats() (models.AnnouncementStats, error) {
	var s models.AnnouncementStats
	err := config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT").Scan(&s.Total)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT WHERE state_annoucement = 'Active' OR state_annoucement IS NULL").Scan(&s.Active)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow("SELECT COUNT(DISTINCT id_announcement) FROM REPORT WHERE id_announcement IS NOT NULL AND status = 'pending'").Scan(&s.Reported)
	return s, err
}

func GetAnnouncementById(id string) (models.Announcement, error) {
	row := config.Conn.QueryRow(announcementSelect+`FROM ANNOUNCEMENT a
		LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		LEFT JOIN USER u ON u.id_user = ua.id_user
		WHERE a.id_announcement = ?`, id)
	return scanAnnouncement(row)
}

func CreateAnnouncement(a models.Announcement) error {
	var idCategory interface{}
	if a.Id_category != "" {
		idCategory = a.Id_category
	}
	_, err := config.Conn.Exec(`
		INSERT INTO ANNOUNCEMENT (
			id_announcement, id_category, title_announcement, address_annoucement,
			city, postal, description_annoucement, availability_date, price,
			request, state_annoucement, type_announcement, condition_announcement
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Id_announcement, idCategory, a.Title_announcement, a.Address_annoucement,
		a.City, a.Postal, a.Description_annoucement, a.Availability_date, a.Price,
		a.Request, a.State_annoucement, a.TypeAnnouncement, a.ConditionAnnouncement,
	)
	return err
}

func CreateUserAnnouncement(a models.Announcement, userID string, photoURLs []string) error {
	if err := CreateAnnouncement(a); err != nil {
		return err
	}
	if _, err := config.Conn.Exec(
		"INSERT INTO USER_ANNOUNCEMENT (id_user, id_announcement) VALUES (?, ?)",
		userID, a.Id_announcement,
	); err != nil {
		return err
	}
	for _, url := range photoURLs {
		if _, err := config.Conn.Exec(
			"INSERT INTO DOCUMENT (id_document, id_user, category, link) VALUES (?, ?, ?, ?)",
			uuid.New().String(), userID, a.Id_announcement, url,
		); err != nil {
			return err
		}
	}
	return nil
}

func ClaimAnnouncement(announcementID, userID string) error {
	_, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET id_buyer = ?, state_annoucement = 'Vendu' WHERE id_announcement = ? AND id_buyer IS NULL",
		userID, announcementID,
	)
	return err
}

func GetUserAcquisitions(userID string) ([]models.Announcement, error) {
	rows, err := config.Conn.Query(announcementSelect+`FROM ANNOUNCEMENT a
		LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		LEFT JOIN USER u ON u.id_user = ua.id_user
		WHERE a.id_buyer = ?
		ORDER BY a.availability_date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Announcement
	for rows.Next() {
		a, err := scanAnnouncement(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}

	for i := range list {
		var code sql.NullString
		var num sql.NullInt64
		_ = config.Conn.QueryRow(`
			SELECT l.access_code, l.locker_number FROM ANNOUNCEMENT_LOCKER al
			JOIN LOCKER l ON l.id_locker = al.id_locker
			WHERE al.id_announcement = ?`, list[i].Id_announcement,
		).Scan(&code, &num)
		list[i].AccessCode = code.String
		list[i].LockerNumber = int(num.Int64)
	}

	return list, nil
}

func UpdateAnnouncement(a models.Announcement) error {
	_, err := config.Conn.Exec(`
		UPDATE ANNOUNCEMENT SET
			id_category = ?, title_announcement = ?, address_annoucement = ?,
			city = ?, postal = ?, description_annoucement = ?,
			availability_date = ?, price = ?, request = ?, state_annoucement = ?,
			type_announcement = ?, condition_announcement = ?
		WHERE id_announcement = ?`,
		a.Id_category, a.Title_announcement, a.Address_annoucement,
		a.City, a.Postal, a.Description_annoucement,
		a.Availability_date, a.Price, a.Request, a.State_annoucement,
		a.TypeAnnouncement, a.ConditionAnnouncement,
		a.Id_announcement,
	)
	return err
}

func GetAnnouncementOwnerID(announcementID string) (string, error) {
	var userID string
	err := config.Conn.QueryRow(
		"SELECT id_user FROM USER_ANNOUNCEMENT WHERE id_announcement = ?",
		announcementID,
	).Scan(&userID)
	return userID, err
}

func DeleteAnnouncement(id string) error {
	rows, err := config.Conn.Query("SELECT link FROM DOCUMENT WHERE category = ? AND link != ''", id)
	if err != nil {
		return err
	}
	var links []string
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err == nil {
			links = append(links, link)
		}
	}
	rows.Close()

	for _, q := range []string{
		"DELETE FROM REPORT WHERE id_announcement = ?",
		"DELETE FROM USER_ANNOUNCEMENT WHERE id_announcement = ?",
		"DELETE FROM ANNOUNCEMENT_LOCKER WHERE id_announcement = ?",
		"DELETE FROM DOCUMENT WHERE category = ?",
	} {
		if _, err := config.Conn.Exec(q, id); err != nil {
			return err
		}
	}
	_, err = config.Conn.Exec("DELETE FROM ANNOUNCEMENT WHERE id_announcement = ?", id)
	if err != nil {
		return err
	}

	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "./uploads"
	}
	for _, link := range links {
		filename := filepath.Base(link)
		os.Remove(filepath.Join(uploadsDir, filename))
	}
	return nil
}
