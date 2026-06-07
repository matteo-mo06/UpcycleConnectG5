package models

type LoginRequest struct {
	Email        string `json:"email"`
	PasswordUser string `json:"password_user"`
}
