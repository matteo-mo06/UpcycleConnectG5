package models

type Announcement struct {
	IdAnnouncement          string  `json:"id"`
	IdCategory              string  `json:"id_category"`
	TitleAnnouncement       string  `json:"title"`
	AddressAnnouncement     string  `json:"address"`
	City                    string  `json:"city"`
	Postal                  string  `json:"postal"`
	DescriptionAnnouncement string  `json:"description"`
	AvailabilityDate        string  `json:"availability_date"`
	Price                   float64 `json:"price"`
	Request                 int     `json:"request"`
	StateAnnouncement       string  `json:"state"`
	RejectionReason         *string `json:"rejection_reason"`
	AuthorId                string  `json:"author_id"`
	AuthorName              string  `json:"author_name"`
	TypeAnnouncement        string  `json:"type"`
	ConditionAnnouncement   string  `json:"condition"`
	FirstPhoto              string  `json:"first_photo"`
	AccessCode              string  `json:"access_code"`
	LockerNumber            int     `json:"locker_number"`
	CreatedAt               string  `json:"created_at"`
	DeletedAt               string  `json:"deleted_at,omitempty"`
	InvoicePath             string  `json:"invoice_path,omitempty"`
	IsFeatured              int     `json:"is_featured"`
	FeaturedUntil           *string `json:"featured_until,omitempty"`
	FeaturedRequestedAt     *string `json:"featured_requested_at,omitempty"`
}

type CreateAnnouncementRequest struct {
	IdCategory  string   `json:"id_category"`
	Title       string   `json:"title"`
	Address     string   `json:"address"`
	City        string   `json:"city"`
	Postal      string   `json:"postal"`
	Description string   `json:"description"`
	AvailDate   string   `json:"availability_date"`
	Price       float64  `json:"price"`
	Type        string   `json:"type"`
	Condition   string   `json:"condition"`
	PhotoURLs   []string `json:"photo_urls"`
}

type AnnouncementStats struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Featured int `json:"featured"`
	Pending  int `json:"pending"`
	Reported int `json:"reported"`
}
