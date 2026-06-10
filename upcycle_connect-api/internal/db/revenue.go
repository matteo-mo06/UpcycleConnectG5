package db

import (
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

func GetRevenueSummary() (models.RevenueSummary, error) {
	var s models.RevenueSummary
	err := config.Conn.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(amount_cents), 0),
			COALESCE(SUM(commission_amount_cents), 0)
		FROM payement
		WHERE status_payement = 'paid'
	`).Scan(&s.TotalTransactions, &s.TotalAmountCents, &s.TotalCommissionCents)
	if err != nil {
		return s, err
	}
	s.CommissionRate, _ = GetCommissionRate()
	return s, nil
}

func GetRevenueTransactions(page, limit int) ([]models.RevenueTransaction, int, error) {
	offset := (page - 1) * limit

	var total int
	if err := config.Conn.QueryRow("SELECT COUNT(*) FROM payement WHERE status_payement = 'paid'").Scan(&total); err != nil {
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
		WHERE p.status_payement = 'paid'
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
