package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllAnnouncements(idCategory string, request int) ([]models.Announcement, error) {
	query := `SELECT id_announcement, id_category, title_announcement, address_annoucement, city, postal,
	                 description_annoucement, availability_date, price, request, state_annoucement
	          FROM ANNOUNCEMENT WHERE 1=1`
	var args []any

	if idCategory != "" {
		query += " AND id_category = ?"
		args = append(args, idCategory)
	}
	if request >= 0 {
		query += " AND request = ?"
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
		)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
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
	_, err := config.Conn.Exec("DELETE FROM REPORT_ANNOUCEMENT WHERE id_announcement = ?", id)
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
