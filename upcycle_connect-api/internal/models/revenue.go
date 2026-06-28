package models

type RevenueSummary struct {
	TotalTransactions      int     `json:"total_transactions"`
	TotalAmountCents       int64   `json:"total_amount_cents"`
	TotalCommissionCents   int64   `json:"total_commission_cents"`
	TotalSubscriptionsCents int64  `json:"total_subscriptions_cents"`
	TotalAdsCents          int64   `json:"total_ads_cents"`
	CommissionRate         float64 `json:"commission_rate"`
}

type RevenueTransaction struct {
	IdPayement            string `json:"id_payement"`
	AnnouncementID        string `json:"announcement_id"`
	AnnouncementTitle     string `json:"announcement_title"`
	BuyerID               string `json:"buyer_id"`
	BuyerName             string `json:"buyer_name"`
	AmountCents           int    `json:"amount_cents"`
	CommissionAmountCents int    `json:"commission_amount_cents"`
	Currency              string `json:"currency"`
	CreatedAt             string `json:"created_at"`
}

type AdPayment struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	UserID     string  `json:"user_id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      string  `json:"email"`
	PriceCents int     `json:"price_cents"`
	PaidAt     *string `json:"paid_at"`
	State      string  `json:"state"`
}
