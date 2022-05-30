package persistent

import (
	"context"
	"os"
	"time"

	"github.com/0w0mewo/mcrc_tgbot/persistent/ent"
	_ "github.com/mattn/go-sqlite3"
)

var DefaultDBConn *ent.Client

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

func open(cfg *Config) (*ent.Client, error) {
	dbconn, err := ent.Open(cfg.Driver, cfg.Dialect)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = dbconn.Schema.Create(ctx)
	if err != nil {
		return nil, err
	}

	return dbconn, nil
}
