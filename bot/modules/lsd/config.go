package lsd

import "github.com/0w0mewo/mcrc_tgbot/config"

type lineStickerDownconf struct {
	config       config.ConfigType

}

func ConfigFrom(c config.ConfigType) *lineStickerDownconf {
	cfg := &lineStickerDownconf{
		config: c,
	}

	cfg.Reload()

	return cfg
}


func (mc *lineStickerDownconf) Reload() {

}
