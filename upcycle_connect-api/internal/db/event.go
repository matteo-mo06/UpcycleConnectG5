package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllEvents() ([]models.Event, error) {
	rows, err := config.Conn.Query(`
		SELECT event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		       event.capacity, event.price_cents, event.id_creator,
		       CONCAT(user.first_name, ' ', user.last_name) AS creator_name,
		       COUNT(inscription.id_user) AS inscription_count
		FROM EVENT event
		LEFT JOIN USER user ON user.id_user = event.id_creator
		LEFT JOIN USER_EVENT_INSCRIPTION inscription ON inscription.id_event = event.id_event
		GROUP BY event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		         event.capacity, event.price_cents, event.id_creator`)
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
			&e.Id_creator,
			&e.CreatorName,
			&e.InscriptionCount,
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
		SELECT id_event, title_event, description_event, date_event, location_event, capacity, price_cents, id_creator
		FROM EVENT WHERE id_event = ?`, id,
	).Scan(
		&e.Id_event,
		&e.Title_event,
		&e.Description_event,
		&e.Date_event,
		&e.Location_event,
		&e.Capacity,
		&e.Price_cents,
		&e.Id_creator,
	)
	return e, err
}

func CreateEvent(e models.Event) error {
	_, err := config.Conn.Exec(`
		INSERT INTO EVENT (
			id_event, title_event, description_event, date_event,
			location_event, capacity, price_cents, id_creator
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		e.Id_event,
		e.Title_event,
		e.Description_event,
		e.Date_event,
		e.Location_event,
		e.Capacity,
		e.Price_cents,
		e.Id_creator,
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
			price_cents = ?,
			id_creator = ?
		WHERE id_event = ?
	`,
		e.Title_event,
		e.Description_event,
		e.Date_event,
		e.Location_event,
		e.Capacity,
		e.Price_cents,
		e.Id_creator,
		e.Id_event,
	)
	return err
}

func DeleteEvent(id string) error {
	_, err := config.Conn.Exec("DELETE FROM USER_EVENT_INSCRIPTION WHERE id_event = ?", id)
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
