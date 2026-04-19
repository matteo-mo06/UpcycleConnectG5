package models

type Document struct {
	Id_document string `json:"id"`
	Id_user     string `json:"id_user"`
	Category    string `json:"category"`
	Link        string `json:"link"`
}

type DepositRequest struct {
	AnnouncementID string `json:"announcement_id"`
	Title          string `json:"announcement_title"`
	Type           string `json:"announcement_type"`
	Status         string `json:"status"`
	LockerNumber   int    `json:"locker_number"`
}

type Locker struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
}

type UserStats struct {
	UpcyclingScore      int `json:"upcycling_score"`
	ActiveAnnouncements int `json:"active_announcements"`
	PendingDeposits     int `json:"pending_deposits"`
	UpcomingEvents      int `json:"upcoming_events"`
}
