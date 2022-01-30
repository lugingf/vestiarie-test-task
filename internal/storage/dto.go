package storage

import (
	"github.com/lugingf/vestiarie-test-task/internal"
)

type Payout struct {
	ID         int64
	UpdateID   string
	SellerID   int64
	Amount     float64
	Currency   internal.Currency
	ItemIDList []int64
	Part       int64
}

type Item struct {
	ID       int64   `db:"id"`
	UpdateID string  `db:"update_id"`
	Name     string  `db:"item_name"`
	Currency string  `db:"currency"`
	Price    float64 `db:"price"`
	SellerID int64   `db:"seller_id"`
}
