package domain

import (
	"github.com/lugingf/vestiarie-test-task/internal"
	"github.com/lugingf/vestiarie-test-task/internal/storage"
	"github.com/pkg/errors"
)

type PayoutService struct {
	Storage *storage.PayoutStorage
}

func NewPayoutService(s *storage.PayoutStorage) *PayoutService {
	return &PayoutService{
		Storage: s,
	}
}

func (p *PayoutService)StorePayouts(items []Item, updateID string) ([]storage.Payout, error) {
	payouts := p.calculatePayouts(items, updateID)

	db := *p.Storage
	err := db.SavePayouts(payouts)
	if err != nil {
		return nil, errors.Wrap(err, "StorePayouts: cannot save payouts")
	}
	return payouts, nil
}

func (p *PayoutService) calculatePayouts(items []Item, updateID string) []storage.Payout {
	payoutsBySellerAndCurrency := make(map[int64]map[internal.Currency]storage.Payout)

	for _, item := range items {
		seller := item.SellerID
		currency := internal.Currency(item.Currency)

		// new seller (and new currency as well)
		if _, ok := payoutsBySellerAndCurrency[seller]; !ok {
			sellersPayoutsByCurrency := make(map[internal.Currency]storage.Payout)
			payout := storage.Payout{
				UpdateID: updateID,
				SellerID: seller,
				Amount:   item.Price,
				Currency: currency,
			}
			sellersPayoutsByCurrency[currency] = payout
			payoutsBySellerAndCurrency[seller] = sellersPayoutsByCurrency
			continue
		}

		// already have the seller, but new currency for them
		if _, ok := payoutsBySellerAndCurrency[seller][currency]; !ok {
			payout := storage.Payout{
				UpdateID: updateID,
				SellerID: seller,
				Amount:   item.Price,
				Currency: currency,
			}
			payoutsBySellerAndCurrency[seller][currency] = payout
			continue
		}

		// already have seller and currency
		payout := payoutsBySellerAndCurrency[seller][currency]
		payout.Amount += item.Price
		payoutsBySellerAndCurrency[seller][currency] = payout
	}

	payouts := make([]storage.Payout, 0)
	for _, currencies := range payoutsBySellerAndCurrency {
		for _, payout := range currencies {
			payouts = append(payouts, payout)
		}
	}

	return payouts
}