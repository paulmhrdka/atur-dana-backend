package responses

type SummaryPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type SummaryTotals struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type CategoryBreakdown struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Type         string  `json:"type"`
	Total        float64 `json:"total"`
	Percentage   float64 `json:"percentage"`
	Count        int64   `json:"count"`
}

type SummaryResponse struct {
	Period     SummaryPeriod       `json:"period"`
	Totals     SummaryTotals       `json:"totals"`
	ByCategory []CategoryBreakdown `json:"by_category"`
}
