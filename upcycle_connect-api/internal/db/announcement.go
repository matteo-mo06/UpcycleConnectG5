package db

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

const announcementSelect = `SELECT a.id_announcement, a.id_category, a.title_announcement, a.address_annoucement, a.city, a.postal,
	a.description_annoucement, a.availability_date, a.price, a.request, a.state_annoucement, a.rejection_reason,
	u.id_user, u.first_name, u.last_name, a.type_announcement, a.condition_announcement,
	(SELECT d.link FROM DOCUMENT d WHERE d.category = a.id_announcement ORDER BY d.id_document LIMIT 1) AS first_photo,
	a.created_at, a.deleted_at `

func scanAnnouncement(row interface{ Scan(...any) error }) (models.Announcement, error) {
	var a models.Announcement
	var idCat, authorId, firstName, lastName, typ, cond, firstPhoto, createdAt, deletedAt, rejReason sql.NullString
	err := row.Scan(
		&a.IdAnnouncement, &idCat, &a.TitleAnnouncement,
		&a.AddressAnnouncement, &a.City, &a.Postal,
		&a.DescriptionAnnouncement, &a.AvailabilityDate, &a.Price,
		&a.Request, &a.StateAnnouncement, &rejReason,
		&authorId, &firstName, &lastName, &typ, &cond, &firstPhoto, &createdAt, &deletedAt,
	)
	if err != nil {
		return a, err
	}
	a.IdCategory = idCat.String
	a.AuthorId = authorId.String
	a.AuthorName = strings.TrimSpace(firstName.String + " " + lastName.String)
	a.TypeAnnouncement = typ.String
	a.ConditionAnnouncement = cond.String
	a.FirstPhoto = firstPhoto.String
	a.CreatedAt = createdAt.String
	a.DeletedAt = deletedAt.String
	if rejReason.Valid && rejReason.String != "" {
		a.RejectionReason = &rejReason.String
	}
	return a, nil
}

func GetAllAnnouncements(idCategory, search, filterType, filterStatus string, request, limit, offset int) ([]models.Announcement, int, error) {
	var where string
	var args []any

	where = `WHERE a.deleted_at IS NULL`
	if filterStatus == "Supprimée" {
		where = `WHERE a.deleted_at IS NOT NULL`
	} else if filterStatus != "" {
		where += " AND a.state_annoucement = ?"
		args = append(args, filterStatus)
	}

	if idCategory != "" {
		where += " AND a.id_category = ?"
		args = append(args, idCategory)
	}
	if request >= 0 {
		where += " AND a.request = ?"
		args = append(args, request)
	}
	if search != "" {
		where += " AND (a.title_announcement LIKE ? OR TRIM(CONCAT(COALESCE(u.first_name,''), ' ', COALESCE(u.last_name,''))) LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	if filterType != "" {
		where += " AND a.type_announcement = ?"
		args = append(args, filterType)
	}

	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := config.Conn.QueryRow(`SELECT COUNT(*) FROM ANNOUNCEMENT a
	          LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
	          LEFT JOIN USER u ON u.id_user = ua.id_user `+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := announcementSelect + `FROM ANNOUNCEMENT a
	          LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
	          LEFT JOIN USER u ON u.id_user = ua.id_user
	          ` + where + ` LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Announcement
	for rows.Next() {
		a, err := scanAnnouncement(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}
	return list, total, nil
}

func GetPublicAnnouncements(search, idCategory, filterType string, limit, offset int) ([]models.Announcement, int, error) {
	where := `WHERE a.state_annoucement = 'Active' AND a.request = 0 AND a.deleted_at IS NULL`
	var args []any

	if search != "" {
		where += " AND (a.title_announcement LIKE ? OR a.description_annoucement LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	if idCategory != "" {
		where += " AND a.id_category = ?"
		args = append(args, idCategory)
	}
	if filterType != "" {
		where += " AND a.type_announcement = ?"
		args = append(args, filterType)
	}

	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT a "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := announcementSelect + `FROM ANNOUNCEMENT a
	          LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
	          LEFT JOIN USER u ON u.id_user = ua.id_user
	          ` + where + ` ORDER BY a.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := config.Conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Announcement
	for rows.Next() {
		a, err := scanAnnouncement(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}
	return list, total, nil
}

func GetUserAnnouncements(userID string) ([]models.Announcement, error) {
	rows, err := config.Conn.Query(announcementSelect+`FROM ANNOUNCEMENT a
		JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		JOIN USER u ON u.id_user = ua.id_user
		WHERE ua.id_user = ? AND a.deleted_at IS NULL
		ORDER BY a.created_at DESC`, userID)
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
			WHERE al.id_announcement = ?`, list[i].IdAnnouncement,
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
	err := config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT WHERE deleted_at IS NULL").Scan(&s.Total)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT WHERE deleted_at IS NULL AND (state_annoucement = 'Active' OR state_annoucement IS NULL)").Scan(&s.Active)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT WHERE deleted_at IS NULL AND state_annoucement = 'En attente'").Scan(&s.Pending)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow("SELECT COUNT(DISTINCT id_announcement) FROM REPORT WHERE id_announcement IS NOT NULL AND status = 'pending'").Scan(&s.Reported)
	return s, err
}

func SetAnnouncementState(id, state string) error {
	_, err := config.Conn.Exec("UPDATE ANNOUNCEMENT SET state_annoucement = ? WHERE id_announcement = ?", state, id)
	return err
}

func RejectAnnouncementWithReason(id, reason string) error {
	_, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET state_annoucement = 'Refusée', rejection_reason = ? WHERE id_announcement = ?",
		reason, id,
	)
	return err
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
	if a.IdCategory != "" {
		idCategory = a.IdCategory
	}
	_, err := config.Conn.Exec(`
		INSERT INTO ANNOUNCEMENT (
			id_announcement, id_category, title_announcement, address_annoucement,
			city, postal, description_annoucement, availability_date, price,
			request, state_annoucement, type_announcement, condition_announcement
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.IdAnnouncement, idCategory, a.TitleAnnouncement, a.AddressAnnouncement,
		a.City, a.Postal, a.DescriptionAnnouncement, a.AvailabilityDate, a.Price,
		a.Request, a.StateAnnouncement, a.TypeAnnouncement, a.ConditionAnnouncement,
	)
	return err
}

func CreateUserAnnouncement(a models.Announcement, userID string, photoURLs []string) error {
	if err := CreateAnnouncement(a); err != nil {
		return err
	}
	if _, err := config.Conn.Exec(
		"INSERT INTO USER_ANNOUNCEMENT (id_user, id_announcement) VALUES (?, ?)",
		userID, a.IdAnnouncement,
	); err != nil {
		return err
	}
	for _, url := range photoURLs {
		if _, err := config.Conn.Exec(
			"INSERT INTO DOCUMENT (id_document, id_user, category, link) VALUES (?, ?, ?, ?)",
			uuid.New().String(), userID, a.IdAnnouncement, url,
		); err != nil {
			return err
		}
	}
	return nil
}

func ClaimAnnouncement(announcementID, userID string) error {
	result, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET id_buyer = ?, state_annoucement = 'Vendu' WHERE id_announcement = ? AND id_buyer IS NULL",
		userID, announcementID,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetUserAcquisitions(userID string) ([]models.Announcement, error) {
	rows, err := config.Conn.Query(announcementSelect+`FROM ANNOUNCEMENT a
		LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		LEFT JOIN USER u ON u.id_user = ua.id_user
		WHERE a.id_buyer = ? AND a.deleted_at IS NULL
		ORDER BY a.created_at DESC`, userID)
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
			WHERE al.id_announcement = ?`, list[i].IdAnnouncement,
		).Scan(&code, &num)
		list[i].AccessCode = code.String
		list[i].LockerNumber = int(num.Int64)
	}

	return list, nil
}

func UpdateAnnouncement(a models.Announcement) error {
	var idCategory interface{}
	if a.IdCategory != "" {
		idCategory = a.IdCategory
	}
	_, err := config.Conn.Exec(`
		UPDATE ANNOUNCEMENT SET
			id_category = ?, title_announcement = ?, address_annoucement = ?,
			city = ?, postal = ?, description_annoucement = ?,
			availability_date = ?, price = ?, request = ?, state_annoucement = ?,
			type_announcement = ?, condition_announcement = ?
		WHERE id_announcement = ?`,
		idCategory, a.TitleAnnouncement, a.AddressAnnouncement,
		a.City, a.Postal, a.DescriptionAnnouncement,
		a.AvailabilityDate, a.Price, a.Request, a.StateAnnouncement,
		a.TypeAnnouncement, a.ConditionAnnouncement,
		a.IdAnnouncement,
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

func MarkAnnouncementSold(announcementID, buyerID string) error {
	tx, err := config.Conn.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"UPDATE ANNOUNCEMENT SET state_annoucement = 'Vendu', id_buyer = ? WHERE id_announcement = ? AND id_buyer IS NULL",
		buyerID, announcementID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteAnnouncement(id string) error {
	_, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET state_annoucement = 'Supprimée', deleted_at = NOW() WHERE id_announcement = ?",
		id,
	)
	return err
}
