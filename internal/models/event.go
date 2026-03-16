package models

type Event struct {
	Id_event          string  `json:"id"`
	Title_event       string  `json:"title"`
	Description_event *string `json:"description"`
	Date_event        *string `json:"date"`
	Location_event    *string `json:"location"`
	Capacity          *int    `json:"capacity"`
	Price_cents       int     `json:"price_cents"`
}
