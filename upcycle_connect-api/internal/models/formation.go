package models

type Formation struct {
	Id_formation          string  `json:"id"`
	Title_formation       string  `json:"title"`
	Description_formation *string `json:"description"`
	Date_formation        *string `json:"date"`
	Location_formation    *string `json:"location"`
	Capacity              *int    `json:"capacity"`
	Level                 string  `json:"level"`
	Duration_hours        *int    `json:"duration_hours"`
	Status                string  `json:"status"`
	RejectionReason       *string `json:"rejection_reason"`
	Id_category           *string `json:"id_category"`
	CategoryName          *string `json:"category_name"`
	Id_creator            *string `json:"id_creator"`
	CreatorName           *string `json:"creator_name"`
	Id_formateur          *string `json:"id_formateur"`
	FormateurName         *string `json:"formateur_name"`
	InscriptionCount      int     `json:"inscription_count"`
	IsRegistered          bool    `json:"is_registered"`
	CreatedAt             string  `json:"created_at"`
}

type CreateFormationRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Date        *string `json:"date"`
	Location    *string `json:"location"`
	Capacity    *int    `json:"capacity"`
	Level       string  `json:"level"`
	DurationH   *int    `json:"duration_hours"`
	IdCategory  *string `json:"id_category"`
}

type UpdateFormationRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Date        *string `json:"date"`
	Location    *string `json:"location"`
	Capacity    *int    `json:"capacity"`
	Level       string  `json:"level"`
	DurationH   *int    `json:"duration_hours"`
	IdCategory  *string `json:"id_category"`
}

type RejectFormationRequest struct {
	Reason *string `json:"reason"`
}
