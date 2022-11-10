package persistent

type Config struct {
	Dialect string
	Driver  string
}

var DefaultDBSetting = &Config{
	Dialect: "file:./data/mcrc_tgbot.db?mode=rwc&cache=shared&_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=8000",
	Driver:  "sqlite3",
}
