package domain

import (
	"github.com/lugingf/vestiarie-test-task/internal"
	"github.com/lugingf/vestiarie-test-task/internal/storage"
	"github.com/lugingf/vestiarie-test-task/resources"
	"github.com/pkg/errors"
)

type PayoutService struct {
	PayoutStorage   *storage.PayoutStorage
	ItemStorage     *storage.ItemStorage
	MaxPayoutAmount float64
}

func NewPayoutService(ps *storage.PayoutStorage, is *storage.ItemStorage) *PayoutService {
	return &PayoutService{
		PayoutStorage: ps,
		ItemStorage:   is,
		// Just for example
		MaxPayoutAmount: 5000,
	}
}

func (p *PayoutService) StorePayouts(items []Item, updateID string) ([]storage.Payout, error) {
	storedItems, err := p.saveItems(items, updateID)
	if err == resources.ErrUpdateIdExists {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "StorePayouts: cannot save payouts")
	}

	payouts := p.calculatePayouts(storedItems, updateID)

	db := *p.PayoutStorage
	err = db.SavePayouts(payouts)
	if err != nil {
		return nil, errors.Wrap(err, "StorePayouts: cannot save payouts")
	}
	return payouts, nil
}

func (p *PayoutService) saveItems(items []Item, updateID string) ([]storage.Item, error) {
	is := *p.ItemStorage
	exists, err := is.CheckUpdateID(updateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check update_id")
	}
	if exists {
		return nil, resources.ErrUpdateIdExists
	}

	itemsToStore := make([]storage.Item, len(items))
	for i, incomingItem := range items {
		item := storage.Item{
			UpdateID: updateID,
			Name:     incomingItem.Name,
			Currency: incomingItem.Currency,
			Price:    incomingItem.Price,
			SellerID: incomingItem.SellerID,
		}
		itemsToStore[i] = item
	}

	err = is.SaveItems(itemsToStore)
	if err != nil {
		return nil, errors.Wrap(err, "cant save items")
	}

	storedItems, err := is.ItemsByUpdateID(updateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stored items")
	}

	return storedItems, nil
}

func (p *PayoutService) calculatePayouts(items []storage.Item, updateID string) []storage.Payout {
	payoutsBySellerAndCurrency := make(map[int64]map[internal.Currency]map[int]storage.Payout)

MAIN:
	for _, item := range items {
		seller := item.SellerID
		currency := internal.Currency(item.Currency)

		// new seller (and new currency as well)
		if _, ok := payoutsBySellerAndCurrency[seller]; !ok {
			sellersPayoutsByCurrency := make(map[internal.Currency]map[int]storage.Payout)
			sellersPayoutsParts := make(map[int]storage.Payout)
			itemIDList := make([]int64, 0)
			itemIDList = append(itemIDList, item.ID)
			payout := storage.Payout{
				UpdateID:   updateID,
				SellerID:   seller,
				Amount:     item.Price,
				Currency:   currency,
				ItemIDList: itemIDList,
				Part:       1,
			}
			sellersPayoutsParts[1] = payout
			sellersPayoutsByCurrency[currency] = sellersPayoutsParts
			payoutsBySellerAndCurrency[seller] = sellersPayoutsByCurrency
			continue
		}

		// already have the seller, but new currency for them
		if _, ok := payoutsBySellerAndCurrency[seller][currency]; !ok {
			sellersPayoutsParts := make(map[int]storage.Payout)
			itemIDList := make([]int64, 0)
			itemIDList = append(itemIDList, item.ID)
			payout := storage.Payout{
				UpdateID:   updateID,
				SellerID:   seller,
				Amount:     item.Price,
				Currency:   currency,
				ItemIDList: itemIDList,
				Part:       1,
			}
			sellersPayoutsParts[1] = payout
			payoutsBySellerAndCurrency[seller][currency] = sellersPayoutsParts
			continue
		}

		// already have seller and currency
		payoutParts := payoutsBySellerAndCurrency[seller][currency]
		for part, payout := range payoutParts {
			newAmount := payout.Amount + item.Price
			if newAmount < p.MaxPayoutAmount {
				payout.Amount = newAmount
				payout.ItemIDList = append(payout.ItemIDList, item.ID)
				payoutsBySellerAndCurrency[seller][currency][part] = payout
				continue MAIN
			}
		}
		itemIDList := make([]int64, 0)
		itemIDList = append(itemIDList, item.ID)
		partsCount := len(payoutsBySellerAndCurrency[seller][currency])
		payout := storage.Payout{
			UpdateID:   updateID,
			SellerID:   seller,
			Amount:     item.Price,
			Currency:   currency,
			ItemIDList: itemIDList,
			Part:       int64(partsCount + 1),
		}

		payoutsBySellerAndCurrency[seller][currency][partsCount+1] = payout
	}

	payouts := make([]storage.Payout, 0)
	for _, currencies := range payoutsBySellerAndCurrency {
		for _, payoutParts := range currencies {
			for _, payout := range payoutParts {
				payouts = append(payouts, payout)
			}
		}
	}

	return payouts
}
