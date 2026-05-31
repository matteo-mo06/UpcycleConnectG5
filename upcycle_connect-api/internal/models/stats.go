package models

type Stats struct {
	UserCount         int `json:"user_count"`
	EventCount        int `json:"event_count"`
	AnnouncementCount int `json:"announcement_count"`
	ReportCount       int `json:"report_count"`
}