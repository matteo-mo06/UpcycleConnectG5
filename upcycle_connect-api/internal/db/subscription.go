package db

import (
	"database/sql"
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	rows, err := config.Conn.Query(`
		SELECT id_plan, name, price_cents, interval_unit, interval_count, is_active, stripe_price_id, stripe_product_id
		FROM subscription_plans ORDER BY price_cents ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.SubscriptionPlan
	for rows.Next() {
		var p models.SubscriptionPlan
		if err := rows.Scan(&p.IdPlan, &p.Name, &p.PriceCents, &p.IntervalUnit, &p.IntervalCount, &p.IsActive, &p.StripePriceID, &p.StripeProductID); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func GetSubscriptionPlanByID(id string) (models.SubscriptionPlan, error) {
	var p models.SubscriptionPlan
	err := config.Conn.QueryRow(`
		SELECT id_plan, name, price_cents, interval_unit, interval_count, is_active, stripe_price_id, stripe_product_id
		FROM subscription_plans WHERE id_plan = ?`, id).
		Scan(&p.IdPlan, &p.Name, &p.PriceCents, &p.IntervalUnit, &p.IntervalCount, &p.IsActive, &p.StripePriceID, &p.StripeProductID)
	return p, err
}

func UpdateSubscriptionPlan(p models.SubscriptionPlan) error {
	_, err := config.Conn.Exec(`
		UPDATE subscription_plans SET name = ?, price_cents = ?, interval_unit = ?, interval_count = ?, is_active = ?, stripe_price_id = ?, stripe_product_id = ?
		WHERE id_plan = ?`,
		p.Name, p.PriceCents, p.IntervalUnit, p.IntervalCount, p.IsActive, p.StripePriceID, p.StripeProductID, p.IdPlan)
	return err
}

func GetUserActiveSubscription(userID string) (*models.Subscription, error) {
	var sub models.Subscription
	var plan models.SubscriptionPlan
	err := config.Conn.QueryRow(`
		SELECT s.id_subscription, s.start_timestamp, s.end_timestamp, s.cancelled, s.cancelled_at,
		       s.stripe_subscription_id, s.stripe_customer_id,
		       sp.id_plan, sp.name, sp.price_cents, sp.interval_unit, sp.interval_count, sp.is_active, sp.stripe_price_id, sp.stripe_product_id
		FROM user_subscription us
		JOIN subscription s ON s.id_subscription = us.id_subscription
		JOIN sub_sub_plan ssp ON ssp.id_subscription = s.id_subscription
		JOIN subscription_plans sp ON sp.id_plan = ssp.id_plan
		WHERE us.id_user = ?
		  AND s.cancelled = 0
		  AND (s.end_timestamp IS NULL OR s.end_timestamp > NOW())
		ORDER BY s.start_timestamp DESC
		LIMIT 1`, userID).
		Scan(&sub.IdSubscription, &sub.StartTimestamp, &sub.EndTimestamp, &sub.Cancelled, &sub.CancelledAt,
			&sub.StripeSubscriptionID, &sub.StripeCustomerID,
			&plan.IdPlan, &plan.Name, &plan.PriceCents, &plan.IntervalUnit, &plan.IntervalCount, &plan.IsActive, &plan.StripePriceID, &plan.StripeProductID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	sub.Plan = &plan
	return &sub, nil
}

func GetUserIDByStripeCustomerID(customerID string) string {
	var userID sql.NullString
	config.Conn.QueryRow(`
		SELECT us.id_user
		FROM user_subscription us
		JOIN subscription s ON s.id_subscription = us.id_subscription
		WHERE s.stripe_customer_id = ?
		ORDER BY s.start_timestamp DESC LIMIT 1`, customerID).Scan(&userID)
	return userID.String
}

func GetUserStripeCustomerID(userID string) string {
	var customerID sql.NullString
	config.Conn.QueryRow(`
		SELECT s.stripe_customer_id
		FROM user_subscription us
		JOIN subscription s ON s.id_subscription = us.id_subscription
		WHERE us.id_user = ? AND s.stripe_customer_id IS NOT NULL
		ORDER BY s.start_timestamp DESC LIMIT 1`, userID).Scan(&customerID)
	return customerID.String
}

func CreateSubscriptionRecord(subID, stripeSubID, stripeCustomerID, planID, userID string) error {
	tx, err := config.Conn.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO subscription (id_subscription, start_timestamp, cancelled, stripe_subscription_id, stripe_customer_id)
		VALUES (?, NOW(), 0, ?, ?)`,
		subID, stripeSubID, stripeCustomerID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`INSERT INTO sub_sub_plan (id_subscription, id_plan) VALUES (?, ?)`, subID, planID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`INSERT INTO user_subscription (id_user, id_subscription) VALUES (?, ?)`, userID, subID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func UpdateSubscriptionByStripeID(stripeSubID string, endTimestamp *string, cancelled bool) error {
	_, err := config.Conn.Exec(`
		UPDATE subscription SET end_timestamp = ?, cancelled = ?,
		cancelled_at = CASE WHEN ? = 1 THEN NOW() ELSE cancelled_at END
		WHERE stripe_subscription_id = ?`,
		endTimestamp, cancelled, cancelled, stripeSubID)
	return err
}

func GetAllSubscriptionsPaginated(page, limit int) ([]models.UserSubscriptionSummary, int, error) {
	var total int
	if err := config.Conn.QueryRow(`
		SELECT COUNT(*)
		FROM user_subscription us
		JOIN subscription s ON s.id_subscription = us.id_subscription`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	rows, err := config.Conn.Query(`
		SELECT u.id_user, u.first_name, u.last_name, u.email,
		       sp.name, sp.price_cents, s.start_timestamp, s.end_timestamp, s.cancelled
		FROM user_subscription us
		JOIN user u ON u.id_user = us.id_user
		JOIN subscription s ON s.id_subscription = us.id_subscription
		JOIN sub_sub_plan ssp ON ssp.id_subscription = s.id_subscription
		JOIN subscription_plans sp ON sp.id_plan = ssp.id_plan
		ORDER BY s.start_timestamp DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.UserSubscriptionSummary
	for rows.Next() {
		var s models.UserSubscriptionSummary
		if err := rows.Scan(&s.UserID, &s.FirstName, &s.LastName, &s.Email, &s.PlanName, &s.PriceCents, &s.StartDate, &s.EndDate, &s.Cancelled); err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}
	return list, total, nil
}

func CreateSubscriptionPlan(p models.SubscriptionPlan) error {
	_, err := config.Conn.Exec(`
		INSERT INTO subscription_plans (id_plan, name, price_cents, interval_unit, interval_count, is_active, stripe_price_id, stripe_product_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.IdPlan, p.Name, p.PriceCents, p.IntervalUnit, p.IntervalCount, p.IsActive, p.StripePriceID, p.StripeProductID)
	return err
}

func SubscriptionExistsByStripeID(stripeSubID string) bool {
	var count int
	config.Conn.QueryRow(`SELECT COUNT(*) FROM subscription WHERE stripe_subscription_id = ?`, stripeSubID).Scan(&count)
	return count > 0
}

func IsUserPremium(userID string) bool {
	sub, err := GetUserActiveSubscription(userID)
	return err == nil && sub != nil
}
