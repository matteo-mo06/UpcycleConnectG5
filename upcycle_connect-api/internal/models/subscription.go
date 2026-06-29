package models

type SubscriptionPlan struct {
	IdPlan          string  `json:"id"`
	Name            string  `json:"name"`
	PriceCents      int     `json:"price_cents"`
	IntervalUnit    string  `json:"interval_unit"`
	IntervalCount   int     `json:"interval_count"`
	IsActive        bool    `json:"is_active"`
	StripePriceID   *string `json:"stripe_price_id"`
	StripeProductID *string `json:"stripe_product_id"`
}

type Subscription struct {
	IdSubscription       string            `json:"id"`
	StartTimestamp       *string           `json:"start_timestamp"`
	EndTimestamp         *string           `json:"end_timestamp"`
	Cancelled            bool              `json:"cancelled"`
	CancelledAt          *string           `json:"cancelled_at"`
	StripeSubscriptionID *string           `json:"stripe_subscription_id"`
	StripeCustomerID     *string           `json:"stripe_customer_id"`
	Plan                 *SubscriptionPlan `json:"plan,omitempty"`
}

type UserSubscriptionSummary struct {
	UserID     string  `json:"user_id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      string  `json:"email"`
	PlanName   string  `json:"plan_name"`
	PriceCents int     `json:"price_cents"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
	Cancelled  bool    `json:"cancelled"`
}

type UpdateSubscriptionPlanRequest struct {
	Name          string `json:"name"`
	PriceCents    int    `json:"price_cents"`
	IntervalUnit  string `json:"interval_unit"`
	IntervalCount int    `json:"interval_count"`
	IsActive      bool   `json:"is_active"`
}
