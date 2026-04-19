package models

type Announcement struct {
	Id_announcement         string  `json:"id"`
	Id_category             string  `json:"id_category"`
	Title_announcement      string  `json:"title"`
	Address_annoucement     string  `json:"address"`
	City                    string  `json:"city"`
	Postal                  string  `json:"postal"`
	Description_annoucement string  `json:"description"`
	Availability_date       string  `json:"availability_date"`
	Price                   float64 `json:"price"`
	Request                 int     `json:"request"`
	State_annoucement       string  `json:"state"`
	AuthorName              string  `json:"author_name"`
	TypeAnnouncement        string  `json:"type"`
	ConditionAnnouncement   string  `json:"condition"`
	FirstPhoto              string  `json:"first_photo"`
	AccessCode              string  `json:"access_code"`
	LockerNumber            int     `json:"locker_number"`
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
	Reported int `json:"reported"`
}
