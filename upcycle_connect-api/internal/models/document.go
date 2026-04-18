package models

type Document struct {
	Id_document string `json:"id"`
	Id_user     string `json:"id_user"`
	Category    string `json:"category"`
	Link        string `json:"link"`
}

type UserStats struct {
	UpcyclingScore      int `json:"upcycling_score"`
	ActiveAnnouncements int `json:"active_announcements"`
	PendingDeposits     int `json:"pending_deposits"`
	UpcomingEvents      int `json:"upcoming_events"`
}
