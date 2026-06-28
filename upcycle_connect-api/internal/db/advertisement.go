package db

import (
	"database/sql"
	"strconv"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func CreateAdvertisement(ad models.Advertisement) error {
	_, err := config.Conn.Exec(`
		INSERT INTO advertisement (id_advertisement, id_user, title, image_url, link_url, state, price_cents)
		VALUES (?, ?, ?, ?, ?, 'pending', ?)`,
		ad.IdAdvertisement, ad.IdUser, ad.Title, ad.ImageURL, ad.LinkURL, ad.PriceCents)
	return err
}

func GetAdvertisementByID(id string) (*models.Advertisement, error) {
	var ad models.Advertisement
	var linkURL, rejectionReason, checkoutSessionID, approvedAt, paidAt sql.NullString
	err := config.Conn.QueryRow(`
		SELECT id_advertisement, id_user, title, image_url, link_url, state, price_cents,
		       rejection_reason, stripe_checkout_session_id, created_at, approved_at, paid_at
		FROM advertisement WHERE id_advertisement = ?`, id).
		Scan(&ad.IdAdvertisement, &ad.IdUser, &ad.Title, &ad.ImageURL, &linkURL, &ad.State, &ad.PriceCents,
			&rejectionReason, &checkoutSessionID, &ad.CreatedAt, &approvedAt, &paidAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if linkURL.Valid {
		ad.LinkURL = &linkURL.String
	}
	if rejectionReason.Valid {
		ad.RejectionReason = &rejectionReason.String
	}
	if checkoutSessionID.Valid {
		ad.StripeCheckoutSessionID = &checkoutSessionID.String
	}
	if approvedAt.Valid {
		ad.ApprovedAt = &approvedAt.String
	}
	if paidAt.Valid {
		ad.PaidAt = &paidAt.String
	}
	return &ad, nil
}

func GetActiveAdvertisements() ([]models.Advertisement, error) {
	rows, err := config.Conn.Query(`
		SELECT id_advertisement, id_user, title, image_url, link_url, price_cents, paid_at
		FROM advertisement WHERE state = 'active' ORDER BY paid_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Advertisement
	for rows.Next() {
		var ad models.Advertisement
		var linkURL, paidAt sql.NullString
		if err := rows.Scan(&ad.IdAdvertisement, &ad.IdUser, &ad.Title, &ad.ImageURL, &linkURL, &ad.PriceCents, &paidAt); err != nil {
			return nil, err
		}
		ad.State = "active"
		if linkURL.Valid {
			ad.LinkURL = &linkURL.String
		}
		if paidAt.Valid {
			ad.PaidAt = &paidAt.String
		}
		list = append(list, ad)
	}
	return list, nil
}

func GetUserAdvertisements(userID string) ([]models.Advertisement, error) {
	rows, err := config.Conn.Query(`
		SELECT id_advertisement, id_user, title, image_url, link_url, state, price_cents,
		       rejection_reason, created_at, approved_at, paid_at
		FROM advertisement WHERE id_user = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Advertisement
	for rows.Next() {
		var ad models.Advertisement
		var linkURL, rejectionReason, approvedAt, paidAt sql.NullString
		if err := rows.Scan(&ad.IdAdvertisement, &ad.IdUser, &ad.Title, &ad.ImageURL, &linkURL, &ad.State, &ad.PriceCents,
			&rejectionReason, &ad.CreatedAt, &approvedAt, &paidAt); err != nil {
			return nil, err
		}
		if linkURL.Valid {
			ad.LinkURL = &linkURL.String
		}
		if rejectionReason.Valid {
			ad.RejectionReason = &rejectionReason.String
		}
		if approvedAt.Valid {
			ad.ApprovedAt = &approvedAt.String
		}
		if paidAt.Valid {
			ad.PaidAt = &paidAt.String
		}
		list = append(list, ad)
	}
	return list, nil
}

func GetAllAdvertisementsPaginated(page, limit int, state, search string) ([]models.Advertisement, int, error) {
	where := "WHERE 1=1"
	args := []any{}
	if state != "" {
		where += " AND a.state = ?"
		args = append(args, state)
	}
	if search != "" {
		s := "%" + search + "%"
		where += " AND (a.title LIKE ? OR u.first_name LIKE ? OR u.last_name LIKE ? OR u.email LIKE ?)"
		args = append(args, s, s, s, s)
	}

	var total int
	countArgs := append([]any{}, args...)
	if err := config.Conn.QueryRow(
		"SELECT COUNT(*) FROM advertisement a JOIN user u ON u.id_user = a.id_user "+where,
		countArgs...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	queryArgs := append(args, limit, offset)
	rows, err := config.Conn.Query(`
		SELECT a.id_advertisement, a.id_user, a.title, a.image_url, a.link_url, a.state, a.price_cents,
		       a.rejection_reason, a.created_at, a.approved_at, a.paid_at,
		       u.first_name, u.last_name, u.email
		FROM advertisement a
		JOIN user u ON u.id_user = a.id_user
		`+where+`
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.Advertisement
	for rows.Next() {
		var ad models.Advertisement
		var linkURL, rejectionReason, approvedAt, paidAt sql.NullString
		if err := rows.Scan(&ad.IdAdvertisement, &ad.IdUser, &ad.Title, &ad.ImageURL, &linkURL, &ad.State, &ad.PriceCents,
			&rejectionReason, &ad.CreatedAt, &approvedAt, &paidAt,
			&ad.UserFirstName, &ad.UserLastName, &ad.UserEmail); err != nil {
			return nil, 0, err
		}
		if linkURL.Valid {
			ad.LinkURL = &linkURL.String
		}
		if rejectionReason.Valid {
			ad.RejectionReason = &rejectionReason.String
		}
		if approvedAt.Valid {
			ad.ApprovedAt = &approvedAt.String
		}
		if paidAt.Valid {
			ad.PaidAt = &paidAt.String
		}
		list = append(list, ad)
	}
	return list, total, nil
}

func ApproveAdvertisement(id string) error {
	_, err := config.Conn.Exec(`
		UPDATE advertisement SET state = 'approved', approved_at = NOW() WHERE id_advertisement = ?`, id)
	return err
}

func RejectAdvertisement(id, reason string) error {
	_, err := config.Conn.Exec(`
		UPDATE advertisement SET state = 'rejected', rejection_reason = ? WHERE id_advertisement = ?`, reason, id)
	return err
}

func ActivateAdvertisement(id, sessionID string) error {
	_, err := config.Conn.Exec(`
		UPDATE advertisement SET state = 'active', paid_at = NOW(), stripe_checkout_session_id = ?
		WHERE id_advertisement = ?`, sessionID, id)
	return err
}

func DeactivateAdvertisement(id string) error {
	_, err := config.Conn.Exec(`
		UPDATE advertisement SET state = 'expired' WHERE id_advertisement = ?`, id)
	return err
}

func ReactivateAdvertisement(id string) error {
	_, err := config.Conn.Exec(`
		UPDATE advertisement SET state = 'active' WHERE id_advertisement = ?`, id)
	return err
}

func SaveAdvertisementCheckoutSession(id, sessionID string) error {
	_, err := config.Conn.Exec(`
		UPDATE advertisement SET stripe_checkout_session_id = ? WHERE id_advertisement = ?`, sessionID, id)
	return err
}

func GetAdvertisementStats() (models.AdvertisementStats, error) {
	var stats models.AdvertisementStats
	rows, err := config.Conn.Query(`SELECT state, COUNT(*) FROM advertisement GROUP BY state`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var state string
		var count int
		if err := rows.Scan(&state, &count); err != nil {
			return stats, err
		}
		stats.Total += count
		switch state {
		case "pending":
			stats.Pending = count
		case "approved":
			stats.Approved = count
		case "active":
			stats.Active = count
		case "rejected":
			stats.Rejected = count
		case "expired":
			stats.Expired = count
		}
	}
	return stats, nil
}

func GetAdvertisementPriceCents() (int, error) {
	var val string
	err := config.Conn.QueryRow(
		"SELECT value_setting FROM platform_settings WHERE key_setting = 'advertisement_price_cents'",
	).Scan(&val)
	if err != nil {
		return 50000, err
	}
	price, err := strconv.Atoi(val)
	if err != nil {
		return 50000, err
	}
	return price, nil
}

func SetAdvertisementPriceCents(price int) error {
	_, err := config.Conn.Exec(
		"UPDATE platform_settings SET value_setting = ? WHERE key_setting = 'advertisement_price_cents'",
		strconv.Itoa(price),
	)
	return err
}
