package models

type User struct {
	IdUser         string  `json:"id"`
	Email          string  `json:"email"`
	PasswordUser   string  `json:"password_user"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	UpcyclingScore int     `json:"upcycling_score"`
	Premium        int     `json:"premium"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
	Profile        string  `json:"profile"`
	AvatarUrl      *string `json:"avatar_url"`
	TutorialDone   bool    `json:"tutorial_done"`
}

type UserResponse struct {
	IdUser         string  `json:"id"`
	Email          string  `json:"email"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	UpcyclingScore int     `json:"upcycling_score"`
	Premium        int     `json:"premium"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
	AvatarUrl      *string `json:"avatar_url"`
	TutorialDone   bool    `json:"tutorial_done"`
}

type UserListItem struct {
	IdUser         string   `json:"id"`
	Email          string   `json:"email"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	UpcyclingScore int      `json:"upcycling_score"`
	Premium        int      `json:"premium"`
	Status         string   `json:"status"`
	CreatedAt      string   `json:"created_at"`
	Roles          []string `json:"roles"`
}

type LoginUserResponse struct {
	IdUser         string   `json:"id"`
	Email          string   `json:"email"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	UpcyclingScore int      `json:"upcycling_score"`
	Premium        int      `json:"premium"`
	Status         string   `json:"status"`
	CreatedAt      string   `json:"created_at"`
	Roles          []string `json:"roles"`
	Permissions    []string `json:"permissions"`
	AvatarUrl      *string  `json:"avatar_url"`
	TutorialDone   bool     `json:"tutorial_done"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		IdUser:         u.IdUser,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		UpcyclingScore: u.UpcyclingScore,
		Premium:        u.Premium,
		Status:         u.Status,
		CreatedAt:      u.CreatedAt,
		AvatarUrl:      u.AvatarUrl,
		TutorialDone:   u.TutorialDone,
	}
}

func (u User) ToLoginResponse(roles []string, permissions []string) LoginUserResponse {
	return LoginUserResponse{
		IdUser:         u.IdUser,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		UpcyclingScore: u.UpcyclingScore,
		Premium:        u.Premium,
		Status:         u.Status,
		CreatedAt:      u.CreatedAt,
		Roles:          roles,
		Permissions:    permissions,
		AvatarUrl:      u.AvatarUrl,
		TutorialDone:   u.TutorialDone,
	}
}
