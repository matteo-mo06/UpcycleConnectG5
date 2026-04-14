package models

type User struct {
	Id_user         string `json:"id"`
	Email           string `json:"email"`
	Password_user   string `json:"password_user"`
	First_name      string `json:"first_name"`
	Last_name       string `json:"last_name"`
	Upcycling_score int    `json:"upcycling_score"`
	Premium         int    `json:"premium"`
	Status          string `json:"status"`
	Created_at      string `json:"created_at"`
}

type UserResponse struct {
	Id_user         string `json:"id"`
	Email           string `json:"email"`
	First_name      string `json:"first_name"`
	Last_name       string `json:"last_name"`
	Upcycling_score int    `json:"upcycling_score"`
	Premium         int    `json:"premium"`
	Status          string `json:"status"`
	Created_at      string `json:"created_at"`
}

// UserListItem est utilisé pour la liste admin : inclut les rôles
type UserListItem struct {
	Id_user         string   `json:"id"`
	Email           string   `json:"email"`
	First_name      string   `json:"first_name"`
	Last_name       string   `json:"last_name"`
	Upcycling_score int      `json:"upcycling_score"`
	Premium         int      `json:"premium"`
	Status          string   `json:"status"`
	Created_at      string   `json:"created_at"`
	Roles           []string `json:"roles"`
}

// LoginUserResponse est utilisé uniquement pour la réponse du login
type LoginUserResponse struct {
	Id_user         string   `json:"id"`
	Email           string   `json:"email"`
	First_name      string   `json:"first_name"`
	Last_name       string   `json:"last_name"`
	Upcycling_score int      `json:"upcycling_score"`
	Premium         int      `json:"premium"`
	Status          string   `json:"status"`
	Created_at      string   `json:"created_at"`
	Roles           []string `json:"roles"`
	Permissions     []string `json:"permissions"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		Id_user:         u.Id_user,
		Email:           u.Email,
		First_name:      u.First_name,
		Last_name:       u.Last_name,
		Upcycling_score: u.Upcycling_score,
		Premium:         u.Premium,
		Status:          u.Status,
		Created_at:      u.Created_at,
	}
}

func (u User) ToLoginResponse(roles []string, permissions []string) LoginUserResponse {
	return LoginUserResponse{
		Id_user:         u.Id_user,
		Email:           u.Email,
		First_name:      u.First_name,
		Last_name:       u.Last_name,
		Upcycling_score: u.Upcycling_score,
		Premium:         u.Premium,
		Status:          u.Status,
		Created_at:      u.Created_at,
		Roles:           roles,
		Permissions:     permissions,
	}
}
