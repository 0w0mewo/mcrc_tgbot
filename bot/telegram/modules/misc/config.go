package misc

import "github.com/0w0mewo/mcrc_tgbot/config"

type miscModConf struct {
	config config.ConfigType
}

func ConfigFrom(c config.ConfigType) *miscModConf {
	cfg := &miscModConf{
		config: c,
	}

	cfg.Reload()

	return cfg
}

func (mc *miscModConf) Reload() {

}
