package models

type Announcement struct {
	Id_announcement         string `json:"id"`
	Id_category             string `json:"id_category"`
	Title_announcement      string `json:"title"`
	Address_annoucement     string `json:"address"`
	City                    string `json:"city"`
	Postal                  string `json:"postal"`
	Description_annoucement string `json:"description"`
	Availability_date       string `json:"availability_date"`
	Price                   int    `json:"price"`
	Request                 int    `json:"request"`
	State_annoucement       string `json:"state"`
}
