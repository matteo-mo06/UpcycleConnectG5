package db

import (
	"database/sql"
	"strconv"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetCommissionRate() (float64, error) {
	var val string
	err := config.Conn.QueryRow(
		"SELECT value_setting FROM platform_settings WHERE key_setting = 'commission_rate'",
	).Scan(&val)
	if err != nil {
		return 5.0, err
	}
	rate, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 5.0, err
	}
	return rate, nil
}

func SetCommissionRate(rate float64) error {
	val := strconv.FormatFloat(rate, 'f', 2, 64)
	_, err := config.Conn.Exec(
		"UPDATE platform_settings SET value_setting = ? WHERE key_setting = 'commission_rate'",
		val,
	)
	return err
}

func StorePayment(id, announcementID, buyerID, currency, providerRef string, amountCents, commissionCents int) error {
	_, err := config.Conn.Exec(`
		INSERT INTO payement (id_payement, amount_cents, currency, provider_payement, provider_ref, status_payement, announcement_id, buyer_id, commission_amount_cents, paid_at)
		VALUES (?, ?, ?, 'stripe', ?, 'paid', ?, ?, ?, NOW())`,
		id, amountCents, currency, providerRef, announcementID, buyerID, commissionCents,
	)
	return err
}

func StoreSubscriptionPayment(invoiceID, stripeSubID, userID, currency string, amountCents int) error {
	_, err := config.Conn.Exec(`
		INSERT INTO payement (id_payement, amount_cents, currency, provider_payement, provider_ref, status_payement, subscription_id, buyer_id, commission_amount_cents, paid_at)
		VALUES (?, ?, ?, 'stripe', ?, 'paid', ?, ?, 0, NOW())
		ON DUPLICATE KEY UPDATE status_payement = status_payement`,
		invoiceID, amountCents, currency, invoiceID, stripeSubID, userID,
	)
	return err
}

func StoreFormationPayment(id, formationID, buyerID, currency, providerRef string, amountCents int) error {
	_, err := config.Conn.Exec(`
		INSERT INTO payement (id_payement, amount_cents, currency, provider_payement, provider_ref, status_payement, formation_id, buyer_id, commission_amount_cents, paid_at)
		VALUES (?, ?, ?, 'stripe', ?, 'paid', ?, ?, 0, NOW())`,
		id, amountCents, currency, providerRef, formationID, buyerID,
	)
	return err
}

func GetRevenueSummary() (models.RevenueSummary, error) {
	var s models.RevenueSummary
	err := config.Conn.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(amount_cents), 0),
			COALESCE(SUM(commission_amount_cents), 0)
		FROM payement
		WHERE status_payement = 'paid'
		AND (announcement_id IS NOT NULL OR formation_id IS NOT NULL)
	`).Scan(&s.TotalTransactions, &s.TotalAmountCents, &s.TotalCommissionCents)
	if err != nil {
		return s, err
	}
	var adCount, subCount int
	config.Conn.QueryRow(`SELECT COUNT(*) FROM advertisement WHERE state IN ('active', 'expired')`).Scan(&adCount)
	config.Conn.QueryRow(`SELECT COUNT(*) FROM user_subscription`).Scan(&subCount)
	s.TotalTransactions += adCount + subCount
	config.Conn.QueryRow(`
		SELECT COALESCE(SUM(amount_cents), 0)
		FROM payement
		WHERE subscription_id IS NOT NULL AND status_payement = 'paid'
	`).Scan(&s.TotalSubscriptionsCents)
	config.Conn.QueryRow(`
		SELECT COALESCE(SUM(price_cents), 0)
		FROM advertisement
		WHERE state IN ('active', 'expired')
	`).Scan(&s.TotalAdsCents)
	config.Conn.QueryRow(`
		SELECT COALESCE(SUM(amount_cents), 0)
		FROM payement
		WHERE formation_id IS NOT NULL AND status_payement = 'paid'
	`).Scan(&s.TotalFormationsCents)
	s.CommissionRate, _ = GetCommissionRate()
	return s, nil
}

func GetAdPaymentsPaginated(page, limit int) ([]models.AdPayment, int, error) {
	var total int
	if err := config.Conn.QueryRow(`
		SELECT COUNT(*) FROM advertisement WHERE state IN ('active', 'expired')
	`).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	rows, err := config.Conn.Query(`
		SELECT a.id_advertisement, a.title, a.id_user,
		       u.first_name, u.last_name, u.email,
		       a.price_cents, a.paid_at, a.state
		FROM advertisement a
		JOIN user u ON u.id_user = a.id_user
		WHERE a.state IN ('active', 'expired')
		ORDER BY a.paid_at DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.AdPayment
	for rows.Next() {
		var p models.AdPayment
		if err := rows.Scan(&p.ID, &p.Title, &p.UserID, &p.FirstName, &p.LastName, &p.Email, &p.PriceCents, &p.PaidAt, &p.State); err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}
	return list, total, nil
}

func GetFormationPaymentsPaginated(page, limit int) ([]models.FormationPayment, int, error) {
	var total int
	if err := config.Conn.QueryRow(`
		SELECT COUNT(*) FROM payement WHERE formation_id IS NOT NULL AND status_payement = 'paid'
	`).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	rows, err := config.Conn.Query(`
		SELECT p.id_payement, p.formation_id, f.title_formation,
		       p.buyer_id, u.first_name, u.last_name, u.email,
		       p.amount_cents, p.paid_at
		FROM payement p
		JOIN formation f ON f.id_formation = p.formation_id
		JOIN user u ON u.id_user = p.buyer_id
		WHERE p.formation_id IS NOT NULL AND p.status_payement = 'paid'
		ORDER BY p.paid_at DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.FormationPayment
	for rows.Next() {
		var fp models.FormationPayment
		var paidAt sql.NullString
		if err := rows.Scan(&fp.ID, &fp.FormationID, &fp.FormationTitle,
			&fp.BuyerID, &fp.FirstName, &fp.LastName, &fp.Email,
			&fp.AmountCents, &paidAt); err != nil {
			return nil, 0, err
		}
		if paidAt.Valid {
			fp.PaidAt = &paidAt.String
		}
		list = append(list, fp)
	}
	return list, total, nil
}

func GetRevenueTransactions(page, limit int) ([]models.RevenueTransaction, int, error) {
	offset := (page - 1) * limit

	var total int
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM payement WHERE status_payement = 'paid' AND announcement_id IS NOT NULL").Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := config.Conn.Query(`
		SELECT
			p.id_payement,
			COALESCE(p.announcement_id, ''),
			COALESCE(a.title_announcement, ''),
			COALESCE(p.buyer_id, ''),
			COALESCE(CONCAT(u.first_name, ' ', u.last_name), ''),
			p.amount_cents,
			p.commission_amount_cents,
			p.currency,
			p.created_at
		FROM payement p
		LEFT JOIN announcement a ON a.id_announcement = p.announcement_id
		LEFT JOIN user u ON u.id_user = p.buyer_id
		WHERE p.status_payement = 'paid' AND p.announcement_id IS NOT NULL
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.RevenueTransaction
	for rows.Next() {
		var t models.RevenueTransaction
		if err := rows.Scan(
			&t.IdPayement, &t.AnnouncementID, &t.AnnouncementTitle,
			&t.BuyerID, &t.BuyerName,
			&t.AmountCents, &t.CommissionAmountCents,
			&t.Currency, &t.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, t)
	}
	return list, total, nil
}
