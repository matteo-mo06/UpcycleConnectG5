package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllEvents() ([]models.Event, error) {
	rows, err := config.Conn.Query(`
		SELECT id_event, title_event, description_event, date_event, location_event, capacity, price_cents
		FROM EVENT`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(
			&e.Id_event,
			&e.Title_event,
			&e.Description_event,
			&e.Date_event,
			&e.Location_event,
			&e.Capacity,
			&e.Price_cents,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func GetEventById(id string) (models.Event, error) {
	var e models.Event
	err := config.Conn.QueryRow(`
		SELECT id_event, title_event, description_event, date_event, location_event, capacity, price_cents
		FROM EVENT WHERE id_event = ?`, id,
	).Scan(
		&e.Id_event,
		&e.Title_event,
		&e.Description_event,
		&e.Date_event,
		&e.Location_event,
		&e.Capacity,
		&e.Price_cents,
	)
	return e, err
}

func CreateEvent(e models.Event) error {
	_, err := config.Conn.Exec(`
		INSERT INTO EVENT (
			id_event,
			title_event,
			description_event,
			date_event,
			location_event,
			capacity,
			price_cents
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		e.Id_event,
		e.Title_event,
		e.Description_event,
		e.Date_event,
		e.Location_event,
		e.Capacity,
		e.Price_cents,
	)
	return err
}

func UpdateEvent(e models.Event) error {
	_, err := config.Conn.Exec(`
		UPDATE EVENT SET
			title_event = ?,
			description_event = ?,
			date_event = ?,
			location_event = ?,
			capacity = ?,
			price_cents = ?
		WHERE id_event = ?
	`,
		e.Title_event,
		e.Description_event,
		e.Date_event,
		e.Location_event,
		e.Capacity,
		e.Price_cents,
		e.Id_event,
	)
	return err
}

func DeleteEvent(id string) error {
	_, err := config.Conn.Exec("DELETE FROM USER_EVENT WHERE id_event = ?", id)
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec("DELETE FROM USER_EVENT_INSCRIPTION WHERE id_event = ?", id)
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec("DELETE FROM PAYEMENT_EVENT WHERE id_event = ?", id)
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec("DELETE FROM EVENT WHERE id_event = ?", id)
	return err
}
