package models

type ScoreAction struct {
	ActionType  string `json:"action_type"`
	Points      int    `json:"points"`
	Description string `json:"description"`
}

type ScoreBreakdownItem struct {
	ActionType  string `json:"action_type"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	Count       int    `json:"count"`
	Subtotal    int    `json:"subtotal"`
}

type ScoreBreakdown struct {
	Score     int                  `json:"score"`
	Breakdown []ScoreBreakdownItem `json:"breakdown"`
}
