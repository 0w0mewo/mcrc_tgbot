package persistent

type Config struct {
	Dialect string
	Driver  string
}

var DefaultDBSetting = &Config{
	Dialect: "file:./data/mcrc_tgbot.db?mode=rwc&cache=shared&_journal_mode=WAL&_fk=1&_sync=1&_cache_size=32768&page_size=32768&temp_store=memory",
	Driver:  "sqlite3",
}
