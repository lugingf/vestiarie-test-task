package domain

type Item struct {
	Name     string  `json:"name,omitempty"`
	Currency string  `json:"currency,omitempty"`
	Price    float64 `json:"price,omitempty"`
	SellerID int64   `json:"seller_id,omitempty"`
}
