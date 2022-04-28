package repeater

import (
	"github.com/0w0mewo/mcrc_tgbot/config"
)

type repeaterConf struct {
	config             config.ConfigType
	randstart, randend int
	qutoes             []string
}

func ConfigFrom(c config.ConfigType) *repeaterConf {
	cfg := &repeaterConf{
		config: c,
	}

	cfg.Reload()

	return cfg
}

func (rc *repeaterConf) GetRandStart() int {
	if re := rc.config["repeater"].(config.ConfigType)["rand"].(config.ConfigType)["start"].(int); re != 0 {
		return re
	}

	return 3
}

func (rc *repeaterConf) GetRandEnd() int {
	if re := rc.config["repeater"].(config.ConfigType)["rand"].(config.ConfigType)["end"].(int); re != 0 {
		return re
	}

	return 5
}

func (c *repeaterConf) GetQutoes() []string {
	q := c.config["qutoes"].([]any)

	ret := make([]string, 0, len(q))

	for _, qut := range q {
		ret = append(ret, qut.(string))
	}

	return ret
}

func (c *repeaterConf) Reload() {
	c.qutoes = c.GetQutoes()
	c.randstart = c.GetRandStart()
	c.randend = c.GetRandEnd()

}
