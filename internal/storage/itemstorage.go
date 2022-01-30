package storage

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type ItemStorage interface {
	SaveItems([]Item) error
	ItemsByUpdateID(string) ([]Item, error)
	CheckUpdateID(string) (bool, error)
}

type ItemStorageSQL struct {
	DB *sql.DB
}

func NewItemStorageSQL(connection *sql.DB) ItemStorage {
	return &ItemStorageSQL{DB: connection}
}

func (p *ItemStorageSQL) SaveItems(items []Item) error {
	if len(items) == 0 {
		return nil
	}
	var values string
	data := make([]interface{}, 0)
	for _, item := range items {
		data = append(data, []interface{}{
			item.UpdateID,
			item.Name,
			item.SellerID,
			item.Price,
			item.Currency,
		}...)
		values = fmt.Sprintf("%v %v", values, "(?, ?, ?, ?, ?),")
	}
	trimmed := strings.Trim(values, ",")
	insertScript := `
		INSERT INTO item (
			 update_id, 
			 item_name,
			 seller_id, 
			 price, 
			 currency)
		VALUES %v
		`
	stmt, err := p.DB.Prepare(fmt.Sprintf(insertScript, trimmed))
	if err != nil {
		return errors.Wrap(err, "SaveItems.Prepare")
	}
	defer stmt.Close()

	_, err = stmt.Exec(data...)
	return errors.Wrap(err, "SaveItems.Exec")
}

func (p *ItemStorageSQL) ItemsByUpdateID(updateID string) ([]Item, error) {
	stmt, err := p.DB.Prepare(`
		SELECT 
		     id,
			 update_id, 
			 item_name,
			 seller_id, 
			 price, 
			 currency
		FROM item
		WHERE update_id = ?
	`)
	if err != nil {
		return nil, errors.Wrapf(err, "ItemsByUpdateID.Prepare updateID: %s", updateID)
	}
	defer stmt.Close()

	items := make([]Item, 0)
	rows, err := stmt.Query(updateID)
	if err == sql.ErrNoRows {
		return items, nil
	}
	if err != nil {
		return items, errors.Wrapf(err, "ItemsByUpdateID.Query updateID: %s", updateID)
	}
	for rows.Next() {
		item := Item{}
		err = rows.Scan(
			&item.ID,
			&item.UpdateID,
			&item.Name,
			&item.SellerID,
			&item.Price,
			&item.Currency,
		)
		if err != nil {
			return items, errors.Wrapf(err, "ItemsByUpdateID.GetRowValue updateID: %s", updateID)
		}
		items = append(items, item)
	}

	return items, nil
}

func (p *ItemStorageSQL) CheckUpdateID(updateID string) (bool, error) {
	stmt, err := p.DB.Prepare(`
		SELECT count(1)
		FROM item
		WHERE update_id = ?
	`)
	if err != nil {
		return false, errors.Wrapf(err, "ItemsByUpdateID.Prepare updateID: %s", updateID)
	}
	defer stmt.Close()

	rows, err := stmt.Query(updateID)
	log.Print(err)
	var count int64
	for rows.Next() {
		rows.Scan(&count)
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}