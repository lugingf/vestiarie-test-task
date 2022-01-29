package resources

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

var Di *ResourceContainer

type ResourceContainer struct {
	AppConfig *Config
	SQLShard  *sql.DB
}

// Init FIXME should be "sync.Once" or something
func Init() *ResourceContainer {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("init app config error: %v", err))
	}

	db, err := initDB(cfg.Storage)
	if err != nil {
		log.Fatal(fmt.Sprintf("init DB error: %v", err))
	}

	err = makeMigrations(db)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to make migrations: %v", err))
	}

	Di = &ResourceContainer{
		AppConfig: cfg,
		SQLShard:  db,
	}

	return Di
}

func initDB(config *DataBaseConfig) (*sql.DB, error) {
	driver := config.Driver
	connectString := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v",
		config.User, config.Password, config.Host, config.Port, config.Name)
	return sql.Open(driver, connectString)
}

func makeMigrations(db *sql.DB) error {
	_, err := db.Exec(`
	  CREATE TABLE IF NOT EXISTS payouts (
		update_id varchar(16),
		seller_id integer ,
		amount integer,
		currency varchar(4),
	    UNIQUE INDEX ux_sel_cur_upd(seller_id, currency, update_id)
	  );
`)
	return errors.Wrap(err, "Can't make migrations")
}
