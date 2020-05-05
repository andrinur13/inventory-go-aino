package entities

type DiscountMultiModel struct {
	Discm_name        string  `json:"discount_name"`
	Discm_start_date  string  `json:"start_date"`
	Discm_end_date    string  `json:"end_date"`
	Discm_destination int     `json:"destination"`
	Discm_value       float32 `json:"value"`
}

func (DiscountMultiModel) TableName() string {
	return "master_discount_multi"
}
