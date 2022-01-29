package storage

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type PayoutStorage interface {
	SavePayouts([]Payout) error
}

type PayoutStorageSQL struct {
	DB *sql.DB
}

func NewPayoutStorageSQL(connection *sql.DB) PayoutStorage {
	return &PayoutStorageSQL{DB: connection}
}

func (p *PayoutStorageSQL) SavePayouts(payouts []Payout) error {
	if len(payouts) == 0 {
		return nil
	}
	var values string
	data := make([]interface{}, 0)
	for _, payout := range payouts {
		data = append(data, []interface{}{
			payout.UpdateID,
			payout.SellerID,
			payout.Amount,
			payout.Currency,
		}...)
		values = fmt.Sprintf("%v %v", values, "(?, ?, ?, ?),")
	}
	trimmed := strings.Trim(values, ",")
	insertScript := `
		INSERT INTO payouts 
		    (
			update_id, 
			seller_id, 
			amount, 
			currency)
		VALUES %v
		ON DUPLICATE KEY UPDATE
			update_id = VALUES(update_id),
			seller_id = VALUES(seller_id),
			amount = VALUES(amount),
			currency = VALUES(currency)
		`
	stmt, err := p.DB.Prepare(fmt.Sprintf(insertScript, trimmed))
	if err != nil {
		return errors.Wrap(err, "payoutsaver.Save.Prepare")
	}
	defer stmt.Close()

	_, err = stmt.Exec(data...)
	return errors.Wrap(err, "payoutsaver.Save.Exec")
}