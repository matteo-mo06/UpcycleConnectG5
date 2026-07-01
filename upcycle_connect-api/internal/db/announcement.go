package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

const announcementSelect = `SELECT a.id_announcement, a.id_category, a.title_announcement, a.address_annoucement, a.city, a.postal,
	a.description_annoucement, a.availability_date, a.price, a.request, a.state_annoucement, a.rejection_reason,
	u.id_user, u.first_name, u.last_name, a.type_announcement, a.condition_announcement,
	(SELECT d.link FROM DOCUMENT d WHERE d.category = a.id_announcement ORDER BY d.id_document LIMIT 1) AS first_photo,
	a.created_at, a.deleted_at, a.is_featured, a.featured_until, a.featured_requested_at `

func scanAnnouncement(row interface{ Scan(...any) error }) (models.Announcement, error) {
	var a models.Announcement
	var idCat, authorId, firstName, lastName, typ, cond, firstPhoto, createdAt, deletedAt, rejReason sql.NullString
	var featuredUntil, featuredRequestedAt sql.NullString
	err := row.Scan(
		&a.IdAnnouncement, &idCat, &a.TitleAnnouncement,
		&a.AddressAnnouncement, &a.City, &a.Postal,
		&a.DescriptionAnnouncement, &a.AvailabilityDate, &a.Price,
		&a.Request, &a.StateAnnouncement, &rejReason,
		&authorId, &firstName, &lastName, &typ, &cond, &firstPhoto, &createdAt, &deletedAt,
		&a.IsFeatured, &featuredUntil, &featuredRequestedAt,
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
	if featuredUntil.Valid {
		a.FeaturedUntil = &featuredUntil.String
	}
	if featuredRequestedAt.Valid {
		a.FeaturedRequestedAt = &featuredRequestedAt.String
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
	_ = ExpireFeatured()
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
	          ` + where + ` ORDER BY a.is_featured DESC, a.created_at DESC LIMIT ? OFFSET ?`
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
	_ = ExpireFeatured()
	var s models.AnnouncementStats
	err := config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT WHERE deleted_at IS NULL").Scan(&s.Total)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT WHERE deleted_at IS NULL AND (state_annoucement = 'Active' OR state_annoucement IS NULL)").Scan(&s.Active)
	if err != nil {
		return s, err
	}
	err = config.Conn.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT WHERE deleted_at IS NULL AND is_featured = 1").Scan(&s.Featured)
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

		var invoicePath sql.NullString
		_ = config.Conn.QueryRow(
			"SELECT invoice_path FROM ANNOUNCEMENT WHERE id_announcement = ?",
			list[i].IdAnnouncement,
		).Scan(&invoicePath)
		list[i].InvoicePath = invoicePath.String
	}

	return list, nil
}

func StoreInvoicePath(announcementID, path string) error {
	_, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET invoice_path = ? WHERE id_announcement = ?",
		path, announcementID,
	)
	return err
}

func GetInvoicePath(announcementID string) (string, error) {
	var path sql.NullString
	err := config.Conn.QueryRow(
		"SELECT invoice_path FROM ANNOUNCEMENT WHERE id_announcement = ?",
		announcementID,
	).Scan(&path)
	return path.String, err
}

func GetAnnouncementBuyerID(announcementID string) (string, error) {
	var buyerID sql.NullString
	err := config.Conn.QueryRow(
		"SELECT id_buyer FROM ANNOUNCEMENT WHERE id_announcement = ?", announcementID,
	).Scan(&buyerID)
	return buyerID.String, err
}

func IsAnnouncementBuyer(announcementID, userID string) (bool, error) {
	var count int
	err := config.Conn.QueryRow(
		"SELECT COUNT(*) FROM ANNOUNCEMENT WHERE id_announcement = ? AND id_buyer = ?",
		announcementID, userID,
	).Scan(&count)
	return count > 0, err
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

var ErrNotOwner = errors.New("not owner")
var ErrNotActive = errors.New("announcement not active")
var ErrMonthlyQuota = errors.New("monthly quota exceeded")
var ErrNoSlot = errors.New("no slot available")
var ErrNotQueued = errors.New("not queued")

func ExpireFeatured() error {
	_, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET is_featured = 0, featured_until = NULL WHERE is_featured = 1 AND featured_until < NOW()",
	)
	return err
}

func GetFeaturedCount() (int, error) {
	var count int
	err := config.Conn.QueryRow(
		"SELECT COUNT(*) FROM ANNOUNCEMENT WHERE is_featured = 1 AND featured_until > NOW()",
	).Scan(&count)
	return count, err
}

func GetNextFeaturedSlot() (time.Time, error) {
	var earliest sql.NullString
	err := config.Conn.QueryRow(
		"SELECT MIN(featured_until) FROM ANNOUNCEMENT WHERE is_featured = 1",
	).Scan(&earliest)
	if err != nil {
		return time.Time{}, err
	}
	if !earliest.Valid {
		return time.Now(), nil
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", earliest.String, time.Local)
	return t, err
}

func RequestFeature(announcementID, userID string) (string, *time.Time, error) {
	ok, err := IsAnnouncementOwner(announcementID, userID)
	if err != nil || !ok {
		return "", nil, ErrNotOwner
	}

	var state string
	err = config.Conn.QueryRow(
		"SELECT state_annoucement FROM ANNOUNCEMENT WHERE id_announcement = ?", announcementID,
	).Scan(&state)
	if err != nil || state != "Active" {
		return "", nil, ErrNotActive
	}

	if !IsUserPremium(userID) {
		return "", nil, errors.New("premium required")
	}

	var used int
	err = config.Conn.QueryRow(`
		SELECT COUNT(*) FROM ANNOUNCEMENT a
		JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		WHERE ua.id_user = ?
		AND a.featured_requested_at IS NOT NULL
		AND YEAR(a.featured_requested_at) = YEAR(NOW())
		AND MONTH(a.featured_requested_at) = MONTH(NOW())
	`, userID).Scan(&used)
	if err != nil {
		return "", nil, err
	}
	if used > 0 {
		return "", nil, ErrMonthlyQuota
	}

	_ = ExpireFeatured()

	count, err := GetFeaturedCount()
	if err != nil {
		return "", nil, err
	}

	if count < 3 {
		featuredUntil := time.Now().Add(7 * 24 * time.Hour)
		_, err = config.Conn.Exec(`
			UPDATE ANNOUNCEMENT SET is_featured = 1, featured_until = ?, featured_requested_at = NOW()
			WHERE id_announcement = ?`,
			featuredUntil.Format("2006-01-02 15:04:05"), announcementID,
		)
		if err != nil {
			return "", nil, err
		}
		_, err = config.Conn.Exec(`
			INSERT INTO featured_slot_log (id_slot, id_announcement, started_at, ends_at)
			VALUES (?, ?, NOW(), ?)`,
			uuid.New().String(), announcementID, featuredUntil.Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			return "", nil, err
		}
		return "active", &featuredUntil, nil
	}

	nextSlot, err := GetNextFeaturedSlot()
	if err != nil {
		return "", nil, err
	}
	estimatedSlot := nextSlot.Add(7 * 24 * time.Hour)
	_, err = config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET featured_requested_at = NOW() WHERE id_announcement = ?",
		announcementID,
	)
	if err != nil {
		return "", nil, err
	}
	return "queued", &estimatedSlot, nil
}

func ConfirmFeatureRequest(announcementID, userID string) (string, *time.Time, error) {
	ok, err := IsAnnouncementOwner(announcementID, userID)
	if err != nil || !ok {
		return "", nil, ErrNotOwner
	}

	var isFeatured int
	var requestedAt sql.NullString
	err = config.Conn.QueryRow(
		"SELECT is_featured, featured_requested_at FROM ANNOUNCEMENT WHERE id_announcement = ?",
		announcementID,
	).Scan(&isFeatured, &requestedAt)
	if err != nil || isFeatured == 1 || !requestedAt.Valid {
		return "", nil, ErrNotQueued
	}

	_ = ExpireFeatured()

	count, err := GetFeaturedCount()
	if err != nil {
		return "", nil, err
	}

	if count >= 3 {
		nextSlot, _ := GetNextFeaturedSlot()
		estimated := nextSlot.Add(7 * 24 * time.Hour)
		return "", &estimated, ErrNoSlot
	}

	featuredUntil := time.Now().Add(7 * 24 * time.Hour)
	_, err = config.Conn.Exec(`
		UPDATE ANNOUNCEMENT SET is_featured = 1, featured_until = ?
		WHERE id_announcement = ?`,
		featuredUntil.Format("2006-01-02 15:04:05"), announcementID,
	)
	if err != nil {
		return "", nil, err
	}
	_, err = config.Conn.Exec(`
		INSERT INTO featured_slot_log (id_slot, id_announcement, started_at, ends_at)
		VALUES (?, ?, NOW(), ?)`,
		uuid.New().String(), announcementID, featuredUntil.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return "", nil, err
	}
	return "active", &featuredUntil, nil
}

func CancelFeatureRequest(announcementID string) error {
	_, err := config.Conn.Exec(`
		UPDATE ANNOUNCEMENT SET is_featured = 0, featured_until = NULL, featured_requested_at = NULL
		WHERE id_announcement = ?`,
		announcementID,
	)
	return err
}

func AdminToggleFeature(announcementID string, activate bool) error {
	if activate {
		featuredUntil := time.Now().Add(7 * 24 * time.Hour)
		_, err := config.Conn.Exec(`
			UPDATE ANNOUNCEMENT SET is_featured = 1, featured_until = ?, featured_requested_at = NOW()
			WHERE id_announcement = ?`,
			featuredUntil.Format("2006-01-02 15:04:05"), announcementID,
		)
		if err != nil {
			return err
		}
		_, err = config.Conn.Exec(`
			INSERT INTO featured_slot_log (id_slot, id_announcement, started_at, ends_at)
			VALUES (?, ?, NOW(), ?)`,
			uuid.New().String(), announcementID, featuredUntil.Format("2006-01-02 15:04:05"),
		)
		return err
	}
	_, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET is_featured = 0, featured_until = NULL, featured_requested_at = NULL WHERE id_announcement = ?",
		announcementID,
	)
	return err
}
