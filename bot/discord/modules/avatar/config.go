package avatar

import "github.com/0w0mewo/mcrc_tgbot/config"

type avatarConf struct {
	config       config.ConfigType
	syncToUserId string
}

func ConfigFrom(c config.ConfigType) *avatarConf {
	cfg := &avatarConf{
		config: c,
	}

	cfg.Reload()

	return cfg
}

func (mc *avatarConf) Reload() {
	mc.syncToUserId = mc.GetSyncToUserId()
}

func (mc *avatarConf) GetSyncToUserId() string {
	if re := mc.config["avatar_sync"].(string); re != "" {
		return re
	}

	return "0"
}
