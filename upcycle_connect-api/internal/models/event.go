package models

type Event struct {
	IdEvent          string  `json:"id"`
	TitleEvent       string  `json:"title"`
	DescriptionEvent *string `json:"description"`
	DateEvent        *string `json:"date"`
	LocationEvent    *string `json:"location"`
	Capacity         *int    `json:"capacity"`
	PriceCents       int     `json:"price_cents"`
	IdCreator        *string `json:"id_creator"`
	CreatorName      *string `json:"creator_name"`
	InscriptionCount int     `json:"inscription_count"`
	IsRegistered     bool    `json:"is_registered"`
	Status           string  `json:"status"`
	RejectionReason  *string `json:"rejection_reason"`
}
