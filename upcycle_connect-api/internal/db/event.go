package db

import (
	"database/sql"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetAllEvents(search, status string, limit, offset int) ([]models.Event, int, error) {
	where := `WHERE event.deleted_at IS NULL`
	var args []any

	if search != "" {
		where += " AND event.title_event LIKE ?"
		args = append(args, "%"+search+"%")
	}
	switch status {
	case "upcoming":
		where += " AND event.date_event > NOW() AND event.status = 'approved'"
	case "past":
		where += " AND event.date_event <= NOW() AND event.status = 'approved'"
	case "pending":
		where += " AND event.status = 'pending'"
	case "approved":
		where += " AND event.status = 'approved'"
	case "rejected":
		where += " AND event.status = 'rejected'"
	}

	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM EVENT event "+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		       event.capacity, event.price_cents, event.id_creator,
		       CONCAT(user.first_name, ' ', user.last_name) AS creator_name,
		       COUNT(inscription.id_user) AS inscription_count,
		       event.status, event.rejection_reason
		FROM EVENT event
		LEFT JOIN USER user ON user.id_user = event.id_creator
		LEFT JOIN USER_EVENT_INSCRIPTION inscription ON inscription.id_event = event.id_event
		`+where+`
		GROUP BY event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		         event.capacity, event.price_cents, event.id_creator, event.status, event.rejection_reason
		ORDER BY event.date_event DESC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Event
	for rows.Next() {
		var e models.Event
		var rejReason sql.NullString
		err := rows.Scan(
			&e.IdEvent, &e.TitleEvent, &e.DescriptionEvent, &e.DateEvent, &e.LocationEvent,
			&e.Capacity, &e.PriceCents, &e.IdCreator, &e.CreatorName, &e.InscriptionCount,
			&e.Status, &rejReason,
		)
		if err != nil {
			return nil, 0, err
		}
		if rejReason.Valid && rejReason.String != "" {
			e.RejectionReason = &rejReason.String
		}
		list = append(list, e)
	}
	return list, total, nil
}

func GetEventById(id string) (models.Event, error) {
	var e models.Event
	var rejReason sql.NullString
	err := config.Conn.QueryRow(`
		SELECT id_event, title_event, description_event, date_event, location_event,
		       capacity, price_cents, id_creator, status, rejection_reason
		FROM EVENT WHERE id_event = ? AND deleted_at IS NULL`, id,
	).Scan(
		&e.IdEvent, &e.TitleEvent, &e.DescriptionEvent, &e.DateEvent, &e.LocationEvent,
		&e.Capacity, &e.PriceCents, &e.IdCreator, &e.Status, &rejReason,
	)
	if rejReason.Valid && rejReason.String != "" {
		e.RejectionReason = &rejReason.String
	}
	return e, err
}

func CreateEvent(e models.Event) error {
	status := e.Status
	if status == "" {
		status = "pending"
	}
	_, err := config.Conn.Exec(`
		INSERT INTO EVENT (
			id_event, title_event, description_event, date_event,
			location_event, capacity, price_cents, id_creator, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		e.IdEvent, e.TitleEvent, e.DescriptionEvent, e.DateEvent,
		e.LocationEvent, e.Capacity, e.PriceCents, e.IdCreator, status,
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
		e.TitleEvent, e.DescriptionEvent, e.DateEvent, e.LocationEvent,
		e.Capacity, e.PriceCents, e.IdCreator, e.IdEvent,
	)
	return err
}

func ApproveEvent(id string) error {
	_, err := config.Conn.Exec(
		"UPDATE EVENT SET status = 'approved', rejection_reason = NULL WHERE id_event = ?", id,
	)
	return err
}

func RejectEvent(id, reason string) error {
	_, err := config.Conn.Exec(
		"UPDATE EVENT SET status = 'rejected', rejection_reason = ? WHERE id_event = ?", reason, id,
	)
	return err
}

func GetPublicEventsForUser(userID, search, status string, limit, offset int) ([]models.Event, int, error) {
	where := `WHERE event.status = 'approved' AND event.deleted_at IS NULL`
	var args []any
	args = append(args, userID)

	if search != "" {
		where += " AND event.title_event LIKE ?"
		args = append(args, "%"+search+"%")
	}
	switch status {
	case "upcoming":
		where += " AND event.date_event > NOW()"
	case "registered":
		where += " AND EXISTS (SELECT 1 FROM USER_EVENT_INSCRIPTION reg WHERE reg.id_user = ? AND reg.id_event = event.id_event)"
		args = append(args, userID)
	}

	countWhere := where
	countArgs := make([]any, len(args)-1)
	copy(countArgs, args[1:])
	var total int
	if err := config.Conn.QueryRow("SELECT COUNT(DISTINCT event.id_event) FROM EVENT event "+countWhere, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := make([]any, len(args))
	copy(queryArgs, args)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := config.Conn.Query(`
		SELECT event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		       event.capacity, event.price_cents, event.id_creator,
		       CONCAT(user.first_name, ' ', user.last_name) AS creator_name,
		       COUNT(inscription.id_user) AS inscription_count,
		       EXISTS(SELECT 1 FROM USER_EVENT_INSCRIPTION i WHERE i.id_user = ? AND i.id_event = event.id_event) AS is_registered,
		       event.status
		FROM EVENT event
		LEFT JOIN USER user ON user.id_user = event.id_creator
		LEFT JOIN USER_EVENT_INSCRIPTION inscription ON inscription.id_event = event.id_event
		`+where+`
		GROUP BY event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		         event.capacity, event.price_cents, event.id_creator, event.status
		ORDER BY event.date_event ASC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(
			&e.IdEvent, &e.TitleEvent, &e.DescriptionEvent, &e.DateEvent, &e.LocationEvent,
			&e.Capacity, &e.PriceCents, &e.IdCreator, &e.CreatorName, &e.InscriptionCount,
			&e.IsRegistered, &e.Status,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, e)
	}
	return list, total, nil
}

func GetMyCreatedEvents(userID string) ([]models.Event, error) {
	rows, err := config.Conn.Query(`
		SELECT event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		       event.capacity, event.price_cents, event.id_creator,
		       CONCAT(user.first_name, ' ', user.last_name) AS creator_name,
		       COUNT(inscription.id_user) AS inscription_count,
		       event.status, event.rejection_reason
		FROM EVENT event
		LEFT JOIN USER user ON user.id_user = event.id_creator
		LEFT JOIN USER_EVENT_INSCRIPTION inscription ON inscription.id_event = event.id_event
		WHERE event.id_creator = ? AND event.deleted_at IS NULL
		GROUP BY event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		         event.capacity, event.price_cents, event.id_creator, event.status, event.rejection_reason
		ORDER BY event.date_event DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Event
	for rows.Next() {
		var e models.Event
		var rejReason sql.NullString
		err := rows.Scan(
			&e.IdEvent, &e.TitleEvent, &e.DescriptionEvent, &e.DateEvent, &e.LocationEvent,
			&e.Capacity, &e.PriceCents, &e.IdCreator, &e.CreatorName, &e.InscriptionCount,
			&e.Status, &rejReason,
		)
		if err != nil {
			return nil, err
		}
		if rejReason.Valid && rejReason.String != "" {
			e.RejectionReason = &rejReason.String
		}
		list = append(list, e)
	}
	return list, nil
}

func GetUserRegisteredEvents(userID string) ([]models.Event, error) {
	rows, err := config.Conn.Query(`
		SELECT event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		       event.capacity, event.price_cents, event.id_creator,
		       CONCAT(user.first_name, ' ', user.last_name) AS creator_name,
		       COUNT(inscription.id_user) AS inscription_count
		FROM EVENT event
		LEFT JOIN USER user ON user.id_user = event.id_creator
		LEFT JOIN USER_EVENT_INSCRIPTION inscription ON inscription.id_event = event.id_event
		WHERE EXISTS (
			SELECT 1 FROM USER_EVENT_INSCRIPTION WHERE id_user = ? AND id_event = event.id_event
		) AND event.deleted_at IS NULL
		GROUP BY event.id_event, event.title_event, event.description_event, event.date_event, event.location_event,
		         event.capacity, event.price_cents, event.id_creator
		ORDER BY event.date_event ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Event
	for rows.Next() {
		var e models.Event
		e.IsRegistered = true
		e.Status = "approved"
		err := rows.Scan(
			&e.IdEvent, &e.TitleEvent, &e.DescriptionEvent, &e.DateEvent, &e.LocationEvent,
			&e.Capacity, &e.PriceCents, &e.IdCreator, &e.CreatorName, &e.InscriptionCount,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func RegisterUserForEvent(userID, eventID string) error {
	_, err := config.Conn.Exec(
		"INSERT IGNORE INTO USER_EVENT_INSCRIPTION (id_user, id_event) VALUES (?, ?)",
		userID, eventID,
	)
	return err
}

func UnregisterUserFromEvent(userID, eventID string) error {
	_, err := config.Conn.Exec(
		"DELETE FROM USER_EVENT_INSCRIPTION WHERE id_user = ? AND id_event = ?",
		userID, eventID,
	)
	return err
}

func DeleteEvent(id string) error {
	_, err := config.Conn.Exec("UPDATE EVENT SET deleted_at = NOW() WHERE id_event = ?", id)
	return err
}
