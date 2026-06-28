package models

type AdvertisementStats struct {
	Total    int `json:"total"`
	Pending  int `json:"pending"`
	Approved int `json:"approved"`
	Active   int `json:"active"`
	Rejected int `json:"rejected"`
	Expired  int `json:"expired"`
}

type Advertisement struct {
	IdAdvertisement         string  `json:"id"`
	IdUser                  string  `json:"id_user"`
	Title                   string  `json:"title"`
	ImageURL                string  `json:"image_url"`
	LinkURL                 *string `json:"link_url"`
	State                   string  `json:"state"`
	PriceCents              int     `json:"price_cents"`
	RejectionReason         *string `json:"rejection_reason"`
	StripeCheckoutSessionID *string `json:"stripe_checkout_session_id,omitempty"`
	CreatedAt               *string `json:"created_at"`
	ApprovedAt              *string `json:"approved_at"`
	PaidAt                  *string `json:"paid_at"`
	UserFirstName           string  `json:"user_first_name,omitempty"`
	UserLastName            string  `json:"user_last_name,omitempty"`
	UserEmail               string  `json:"user_email,omitempty"`
}
