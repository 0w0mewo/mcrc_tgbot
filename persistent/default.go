package persistent

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DefaultDBConn *sql.DB

func init() {
	var err error

	dbsetting := os.Getenv("DB_SETTING")
	if dbsetting != "" {
		DefaultDBSetting.Dialect = dbsetting
	}

	dbdriver := os.Getenv("DB_DRIVER")
	if dbdriver != "" {
		DefaultDBSetting.Driver = dbdriver
	}

	DefaultDBConn, err = open(DefaultDBSetting)
	if err != nil {
		panic(err)
	}
}

func open(cfg *Config) (*sql.DB, error) {
	DB, err := sql.Open(cfg.Driver, cfg.Dialect)
	if err != nil {
		panic(err)
	}

	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(20)

	return DB, nil
}
