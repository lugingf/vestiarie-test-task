package storage

import (
	"github.com/google/uuid"
	"github.com/lugingf/vestiarie-test-task/internal"
)

type Payout struct {
	UpdateID uuid.UUID
	SellerID int64
	Amount   float64
	Currency internal.Currency
}
