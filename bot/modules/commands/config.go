package commands

import (
	"github.com/0w0mewo/mcrc_tgbot/config"
)

type commanderConf struct {
	config config.ConfigType
	qutoes []string
}

func ConfigFrom(c config.ConfigType) *commanderConf {
	cfg := &commanderConf{
		config: c,
	}

	cfg.Reload()

	return cfg
}

func (c *commanderConf) GetQutoes() []string {
	q := c.config["qutoes"].([]any)

	ret := make([]string, 0, len(q))

	for _, qut := range q {
		ret = append(ret, qut.(string))
	}

	return ret
}

func (c *commanderConf) Reload() {
	c.qutoes = c.GetQutoes()
}
