package models

type Formation struct {
	IdFormation          string  `json:"id"`
	TitleFormation       string  `json:"title"`
	DescriptionFormation *string `json:"description"`
	DateFormation        *string `json:"date"`
	LocationFormation    *string `json:"location"`
	Capacity             *int    `json:"capacity"`
	Level                string  `json:"level"`
	DurationHours        *int    `json:"duration_hours"`
	Status               string  `json:"status"`
	RejectionReason      *string `json:"rejection_reason"`
	IdCategory           *string `json:"id_category"`
	CategoryName         *string `json:"category_name"`
	IdCreator            *string `json:"id_creator"`
	CreatorName          *string `json:"creator_name"`
	IdFormateur          *string `json:"id_formateur"`
	FormateurName        *string `json:"formateur_name"`
	Syllabus             *string `json:"syllabus"`
	Prerequisites        *string `json:"prerequisites"`
	Objectives           *string `json:"objectives"`
	Price            *float64 `json:"price"`
	InscriptionCount int      `json:"inscription_count"`
	IsRegistered     bool     `json:"is_registered"`
	CreatedAt        string   `json:"created_at"`
}

type CreateFormationRequest struct {
	Title         string   `json:"title"`
	Description   *string  `json:"description"`
	Date          *string  `json:"date"`
	Location      *string  `json:"location"`
	Capacity      *int     `json:"capacity"`
	Level         string   `json:"level"`
	DurationH     *int     `json:"duration_hours"`
	IdCategory    *string  `json:"id_category"`
	Price         *float64 `json:"price"`
	Syllabus      *string  `json:"syllabus"`
	Prerequisites *string  `json:"prerequisites"`
	Objectives    *string  `json:"objectives"`
}

type UpdateFormationRequest struct {
	Title         string   `json:"title"`
	Description   *string  `json:"description"`
	Date          *string  `json:"date"`
	Location      *string  `json:"location"`
	Capacity      *int     `json:"capacity"`
	Level         string   `json:"level"`
	DurationH     *int     `json:"duration_hours"`
	IdCategory    *string  `json:"id_category"`
	Price         *float64 `json:"price"`
	Syllabus      *string  `json:"syllabus"`
	Prerequisites *string  `json:"prerequisites"`
	Objectives    *string  `json:"objectives"`
}

type RejectFormationRequest struct {
	Reason *string `json:"reason"`
}
