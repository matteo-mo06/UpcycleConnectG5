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

type ProjectMaterial struct {
	IdMaterial string  `json:"id"`
	IdProject  string  `json:"id_project"`
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	Unit       *string `json:"unit"`
}

type CreateMaterialRequest struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     *string `json:"unit"`
}

type UpdateMaterialRequest struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     *string `json:"unit"`
}

type ProjectStep struct {
	IdStep      string  `json:"id"`
	IdProject   string  `json:"id_project"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
	StepOrder   int     `json:"step_order"`
}

type CreateStepRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StepOrder   int     `json:"step_order"`
}

type UpdateStepRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StepOrder   int     `json:"step_order"`
}

type UpdateStepStatusRequest struct {
	Status string `json:"status"`
}
