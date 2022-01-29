package storage

import (
	"github.com/lugingf/vestiarie-test-task/internal"
)

type Payout struct {
	UpdateID string
	SellerID int64
	Amount   float64
	Currency internal.Currency
}
