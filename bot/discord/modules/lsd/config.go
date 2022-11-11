package lsd

import "github.com/0w0mewo/mcrc_tgbot/config"

type lsdConf struct {
	config       config.ConfigType

}

func ConfigFrom(c config.ConfigType) *lsdConf {
	cfg := &lsdConf{
		config: c,
	}

	cfg.Reload()

	return cfg
}


func (mc *lsdConf) Reload() {

}
