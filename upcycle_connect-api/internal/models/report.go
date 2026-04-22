package models

type Report struct {
	IdReport         string `json:"id_report"`
	IdReporter       string `json:"id_reporter"`
	IdReportedUser   string `json:"id_reported_user"`
	IdAnnouncement   string `json:"id_announcement"`
	IdTopic          string `json:"id_topic"`
	IdPost           string `json:"id_post"`
	Reason           string `json:"reason"`
	Status           string `json:"status"`
	ActionTaken      string `json:"action_taken"`
	ResolvedBy       string `json:"resolved_by"`
	ResolvedAt       string `json:"resolved_at"`
	CreatedAt        string `json:"created_at"`
	ReporterName     string `json:"reporter_name"`
	ReporterEmail    string `json:"reporter_email"`
	ReportedUserName string `json:"reported_user_name"`
	ContentType      string `json:"content_type"`
	ContentTitle     string `json:"content_title"`
}

type Sanction struct {
	IdSanction   string `json:"id_sanction"`
	IdUser       string `json:"id_user"`
	IdAdmin      string `json:"id_admin"`
	IdReport     string `json:"id_report"`
	Type         string `json:"type"`
	Reason       string `json:"reason"`
	DurationDays int    `json:"duration_days"`
	ExpiresAt    string `json:"expires_at"`
	CreatedAt    string `json:"created_at"`
	AdminName    string `json:"admin_name"`
}

type UserHistorySummary struct {
	IdUser        string `json:"id_user"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Status        string `json:"status"`
	ReportCount   int    `json:"report_count"`
	SanctionCount int    `json:"sanction_count"`
}

type UserHistory struct {
	Status          string     `json:"status"`
	Sanctions       []Sanction `json:"sanctions"`
	ReportsReceived []Report   `json:"reports_received"`
}
