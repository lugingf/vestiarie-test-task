package storage

import (
	"database/sql"
	"encoding/json"
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
		itemListb, err := json.Marshal(payout.ItemIDList)
		if err != nil {
			return errors.Wrap(err, "SavePayouts.Marshal Item List")
		}
		data = append(data, []interface{}{
			payout.UpdateID,
			payout.SellerID,
			payout.Amount,
			payout.Currency,
			string(itemListb),
		}...)
		values = fmt.Sprintf("%v %v", values, "(?, ?, ?, ?, ?),")
	}
	trimmed := strings.Trim(values, ",")
	insertScript := `
		INSERT INTO payout 
		    (
			update_id, 
			seller_id, 
			amount, 
			currency,
			item_id_list)
		VALUES %v
		ON DUPLICATE KEY UPDATE
			update_id = VALUES(update_id)
		`
	stmt, err := p.DB.Prepare(fmt.Sprintf(insertScript, trimmed))
	if err != nil {
		return errors.Wrap(err, "SavePayouts.Prepare")
	}
	defer stmt.Close()

	_, err = stmt.Exec(data...)
	return errors.Wrap(err, "SavePayouts.Exec")
}
