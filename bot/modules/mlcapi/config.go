package mlcapi

import "github.com/0w0mewo/mcrc_tgbot/config"

type mlcApiConf struct {
	config config.ConfigType
	lockto []string
}

func ConfigFrom(c config.ConfigType) *mlcApiConf {
	cfg := &mlcApiConf{
		config: c,
	}

	cfg.Reload()

	return cfg
}

func (mc *mlcApiConf) GetMasters() []string {
	masters := mc.config["mlcapi"].(config.ConfigType)["lockto"].([]any)

	ret := make([]string, 0, len(masters))

	for _, m := range masters {
		ret = append(ret, m.(string))
	}

	return ret

}

func (mc *mlcApiConf) Reload() {
	mc.lockto = mc.GetMasters()

}
