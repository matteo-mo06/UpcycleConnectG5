package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllAnnouncements(idCategory string, request int) ([]models.Announcement, error) {
	query := `SELECT a.id_announcement, a.id_category, a.title_announcement, a.address_annoucement, a.city, a.postal,
	                 a.description_annoucement, a.availability_date, a.price, a.request, a.state_annoucement,
	                 COALESCE(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), '')
	          FROM ANNOUNCEMENT a
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
		var a models.Announcement
		err := rows.Scan(
			&a.Id_announcement,
			&a.Id_category,
			&a.Title_announcement,
			&a.Address_annoucement,
			&a.City,
			&a.Postal,
			&a.Description_annoucement,
			&a.Availability_date,
			&a.Price,
			&a.Request,
			&a.State_annoucement,
			&a.AuthorName,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
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
	var a models.Announcement
	err := config.Conn.QueryRow(`
		SELECT id_announcement, id_category, title_announcement, address_annoucement, city, postal,
		       description_annoucement, availability_date, price, request, state_annoucement
		FROM ANNOUNCEMENT WHERE id_announcement = ?`, id,
	).Scan(
		&a.Id_announcement,
		&a.Id_category,
		&a.Title_announcement,
		&a.Address_annoucement,
		&a.City,
		&a.Postal,
		&a.Description_annoucement,
		&a.Availability_date,
		&a.Price,
		&a.Request,
		&a.State_annoucement,
	)
	return a, err
}

func CreateAnnouncement(a models.Announcement) error {
	_, err := config.Conn.Exec(`
		INSERT INTO ANNOUNCEMENT (
			id_announcement,
			id_category,
			title_announcement,
			address_annoucement,
			city,
			postal,
			description_annoucement,
			availability_date,
			price,
			request,
			state_annoucement
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		a.Id_announcement,
		a.Id_category,
		a.Title_announcement,
		a.Address_annoucement,
		a.City,
		a.Postal,
		a.Description_annoucement,
		a.Availability_date,
		a.Price,
		a.Request,
		a.State_annoucement,
	)
	return err
}

func UpdateAnnouncement(a models.Announcement) error {
	_, err := config.Conn.Exec(`
		UPDATE ANNOUNCEMENT SET
			id_category = ?,
			title_announcement = ?,
			address_annoucement = ?,
			city = ?,
			postal = ?,
			description_annoucement = ?,
			availability_date = ?,
			price = ?,
			request = ?,
			state_annoucement = ?
		WHERE id_announcement = ?
	`,
		a.Id_category,
		a.Title_announcement,
		a.Address_annoucement,
		a.City,
		a.Postal,
		a.Description_annoucement,
		a.Availability_date,
		a.Price,
		a.Request,
		a.State_annoucement,
		a.Id_announcement,
	)
	return err
}

func DeleteAnnouncement(id string) error {
	_, err := config.Conn.Exec("DELETE FROM REPORT WHERE id_announcement = ?", id)
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec("DELETE FROM USER_ANNOUNCEMENT WHERE id_announcement = ?", id)
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec("DELETE FROM ANNOUNCEMENT_LOCKER WHERE id_announcement = ?", id)
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec("DELETE FROM ANNOUNCEMENT WHERE id_announcement = ?", id)
	return err
}
