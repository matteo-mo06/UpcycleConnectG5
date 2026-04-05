package models

type User struct {
	Id_user         string `json:"id"`
	Email           string `json:"email"`
	Password_user   string `json:"password_user"`
	First_name      string `json:"first_name"`
	Last_name       string `json:"last_name"`
	Upcycling_score int    `json:"upcycling_score"`
	Premium         int    `json:"premium"`
}

type UserResponse struct {
	Id_user         string `json:"id"`
	Email           string `json:"email"`
	First_name      string `json:"first_name"`
	Last_name       string `json:"last_name"`
	Upcycling_score int    `json:"upcycling_score"`
	Premium         int    `json:"premium"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		Id_user:         u.Id_user,
		Email:           u.Email,
		First_name:      u.First_name,
		Last_name:       u.Last_name,
		Upcycling_score: u.Upcycling_score,
		Premium:         u.Premium,
	}
}
