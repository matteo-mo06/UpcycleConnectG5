package models

type ProfessionalRequest struct {
	Id_request  string  `json:"id"`
	Id_user     string  `json:"user_id"`
	Status      string  `json:"status"` // pending | approved | rejected
	CreatedAt   string  `json:"created_at"`
	ProcessedAt *string `json:"processed_at"` // NULL tant que non traitée donc je mets une adresse pour avoir la possibilité du NULL en bdd
}

type ProfessionalRequestDetail struct {
	ProfessionalRequest
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}
