package models

type ProfessionalRequest struct {
	IdRequest   string  `json:"id"`
	IdUser      string  `json:"user_id"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	ProcessedAt *string `json:"processed_at"`
}

type ProfessionalRequestDetail struct {
	ProfessionalRequest
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
