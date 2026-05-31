package models

type LoginRequest struct {
	Email         string `json:"email"`
	Password_user string `json:"password_user"`
}