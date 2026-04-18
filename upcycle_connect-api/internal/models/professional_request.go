package models

type ProfessionalRequest struct {
	Id_request  string  `json:"id"`
	Id_user     string  `json:"user_id"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	ProcessedAt *string `json:"processed_at"`
}

type ProfessionalRequestDetail struct {
	ProfessionalRequest
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}
