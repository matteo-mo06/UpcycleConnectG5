package models

type Advice struct {
	IdAdvice     string  `json:"id"`
	Title        string  `json:"title"`
	IdCategory   *string `json:"id_category"`
	CategoryName *string `json:"category_name"`
	Description  *string `json:"description"`
	DateAdvice   *string `json:"date_advice"`
	Priority     int     `json:"priority"`
	IsExpired    bool    `json:"is_expired"`
	IdCreator    *string `json:"id_creator"`
	CreatorName  *string `json:"creator_name"`
	CreatedAt    string  `json:"created_at"`
}

type CreateAdviceRequest struct {
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	IdCategory  *string  `json:"id_category"`
	DateAdvice  *string  `json:"date_advice"`
	Priority    int      `json:"priority"`
	PhotoURLs   []string `json:"photo_urls"`
}

type UpdateAdviceRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	IdCategory  *string `json:"id_category"`
	DateAdvice  *string `json:"date_advice"`
	Priority    int     `json:"priority"`
}

type DeletedAdvice struct {
	IdAdvice     string     `json:"id"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	IdCategory   *string    `json:"id_category"`
	CategoryName *string    `json:"category_name"`
	Priority     int        `json:"priority"`
	DateAdvice   *string    `json:"date_advice"`
	IdCreator    *string    `json:"id_creator"`
	CreatorName  *string    `json:"creator_name"`
	DeletedAt    string     `json:"deleted_at"`
	Documents    []Document `json:"documents"`
}
