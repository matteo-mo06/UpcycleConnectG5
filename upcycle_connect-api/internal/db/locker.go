package db

import (
	"database/sql"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

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

func SetDepositRequest(announcementID string, value int) error {
	_, err := config.Conn.Exec(
		"UPDATE ANNOUNCEMENT SET request = ? WHERE id_announcement = ?",
		value, announcementID,
	)
	return err
}
