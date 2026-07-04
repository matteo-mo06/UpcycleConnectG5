package models

type Participant struct {
	Id        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	AvatarUrl *string `json:"avatar_url"`
}
