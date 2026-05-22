package db

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllLockers() ([]models.Locker, error) {
	rows, err := config.Conn.Query(`
		SELECT l.id_locker, l.locker_number, l.access_code,
		       al.id_announcement IS NOT NULL AS occupied,
		       a.title_announcement,
		       CONCAT_WS(' ', seller.first_name, seller.last_name),
		       CONCAT_WS(' ', buyer.first_name, buyer.last_name)
		FROM LOCKER l
		LEFT JOIN ANNOUNCEMENT_LOCKER al ON al.id_locker = l.id_locker
		LEFT JOIN ANNOUNCEMENT a ON a.id_announcement = al.id_announcement
		LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		LEFT JOIN USER seller ON seller.id_user = ua.id_user
		LEFT JOIN USER buyer ON buyer.id_user = a.id_buyer
		ORDER BY l.locker_number ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Locker
	for rows.Next() {
		var l models.Locker
		var accessCode, title, seller, buyer sql.NullString
		if err := rows.Scan(&l.ID, &l.Number, &accessCode, &l.Occupied, &title, &seller, &buyer); err != nil {
			return nil, err
		}
		l.AccessCode = accessCode.String
		l.AnnouncementTitle = title.String
		l.SellerName = seller.String
		l.BuyerName = buyer.String
		list = append(list, l)
	}
	return list, nil
}

func CreateLocker() (models.Locker, error) {
	rows, err := config.Conn.Query("SELECT locker_number FROM LOCKER ORDER BY locker_number ASC")
	if err != nil {
		return models.Locker{}, err
	}
	defer rows.Close()

	next := 1
	for rows.Next() {
		var n int
		if err := rows.Scan(&n); err != nil {
			return models.Locker{}, err
		}
		if n == next {
			next++
		} else {
			break
		}
	}

	l := models.Locker{
		ID:     uuid.New().String(),
		Number: next,
	}
	_, err = config.Conn.Exec(
		"INSERT INTO LOCKER (id_locker, locker_number) VALUES (?, ?)",
		l.ID, l.Number,
	)
	return l, err
}

func DeleteLocker(id string) error {
	var count int
	if err := config.Conn.QueryRow(
		"SELECT COUNT(id_locker) FROM ANNOUNCEMENT_LOCKER WHERE id_locker = ?", id,
	).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("locker occupé")
	}
	_, err := config.Conn.Exec("DELETE FROM LOCKER WHERE id_locker = ?", id)
	return err
}

func GetPendingDepositRequests() ([]models.Announcement, error) {
	rows, err := config.Conn.Query(announcementSelect + `FROM ANNOUNCEMENT a
		LEFT JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		LEFT JOIN USER u ON u.id_user = ua.id_user
		LEFT JOIN ANNOUNCEMENT_LOCKER al ON al.id_announcement = a.id_announcement
		WHERE a.request = 1 AND al.id_locker IS NULL
		ORDER BY a.availability_date ASC`)
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

func GetUserDepositRequests(userID string) ([]models.DepositRequest, error) {
	rows, err := config.Conn.Query(`
		SELECT
			a.id_announcement    AS announcement_id,
			a.title_announcement AS announcement_title,
			a.type_announcement  AS announcement_type,
			l.locker_number      AS locker_number
		FROM ANNOUNCEMENT a
		JOIN USER_ANNOUNCEMENT ua ON ua.id_announcement = a.id_announcement
		LEFT JOIN ANNOUNCEMENT_LOCKER al ON al.id_announcement = a.id_announcement
		LEFT JOIN LOCKER l ON l.id_locker = al.id_locker
		WHERE ua.id_user = ? AND a.request = 1
		ORDER BY a.id_announcement`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DepositRequest
	for rows.Next() {
		var d models.DepositRequest
		var lockerNumber sql.NullInt64
		if err := rows.Scan(&d.AnnouncementID, &d.Title, &d.Type, &lockerNumber); err != nil {
			return nil, err
		}
		d.LockerNumber = int(lockerNumber.Int64)
		if d.LockerNumber > 0 {
			d.Status = "validated"
		} else {
			d.Status = "pending"
		}
		list = append(list, d)
	}
	return list, nil
}

func GetAvailableLocker() (models.Locker, error) {
	var l models.Locker
	err := config.Conn.QueryRow(`
		SELECT id_locker, locker_number FROM LOCKER
		WHERE id_locker NOT IN (SELECT id_locker FROM ANNOUNCEMENT_LOCKER)
		ORDER BY locker_number ASC
		LIMIT 1`).Scan(&l.ID, &l.Number)
	return l, err
}

func AssignLocker(announcementID, lockerID, accessCode string) error {
	_, err := config.Conn.Exec(
		"INSERT INTO ANNOUNCEMENT_LOCKER (id_announcement, id_locker) VALUES (?, ?)",
		announcementID, lockerID,
	)
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec(
		"UPDATE LOCKER SET access_code = ? WHERE id_locker = ?",
		accessCode, lockerID,
	)
	return err
}

func UpdateLockerAccessCode(id, code string) error {
	var count int
	if err := config.Conn.QueryRow(
		"SELECT COUNT(id_locker) FROM ANNOUNCEMENT_LOCKER WHERE id_locker = ?", id,
	).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return errors.New("locker libre")
	}
	_, err := config.Conn.Exec(
		"UPDATE LOCKER SET access_code = ? WHERE id_locker = ?",
		code, id,
	)
	return err
}

func SetDepositRequest(announcementID string, value int) error {
	_, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET request = ? WHERE id_announcement = ?",
		value, announcementID,
	)
	return err
}
