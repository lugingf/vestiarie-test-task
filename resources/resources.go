package resources

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
