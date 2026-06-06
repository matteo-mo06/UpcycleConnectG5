package models

type Project struct {
	IdProject          string  `json:"id"`
	TitleProject       string  `json:"title"`
	DescriptionProject *string `json:"description"`
	StartDateProject   *string `json:"start_date"`
	EndDate            *string `json:"end_date"`
	LocationProject    *string `json:"location"`
	Capacity           *int    `json:"capacity"`
	Status             string  `json:"status"`
	RejectionReason    *string `json:"rejection_reason"`
	IdCreator          *string `json:"id_creator"`
	CreatorName        *string `json:"creator_name"`
	MembersCount       int     `json:"members_count"`
	IsRegistered       bool    `json:"is_registered"`
	CreatedAt          string  `json:"created_at"`
	UpVotes            int     `json:"up_votes"`
	DownVotes          int     `json:"down_votes"`
	UserVote           int     `json:"user_vote"`
}

type CreateProjectRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Location    *string `json:"location"`
	Capacity    *int    `json:"capacity"`
}

type UpdateProjectAdminRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Location    *string `json:"location"`
	Capacity    *int    `json:"capacity"`
	Status      string  `json:"status"`
}

type RejectProjectRequest struct {
	Reason *string `json:"reason"`
}

type ProjectStats struct {
	Total      int `json:"total"`
	Pending    int `json:"pending"`
	Open       int `json:"open"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Rejected   int `json:"rejected"`
}
