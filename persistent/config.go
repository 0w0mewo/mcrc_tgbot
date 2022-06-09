package persistent

type Config struct {
	Dialect string
	Driver  string
}

var DefaultDBSetting = &Config{
	Dialect: "file:./data/mcrc_tgbot.db?mode=rwc&cache=shared&_journal_mode=WAL",
	Driver:  "sqlite3",
}
