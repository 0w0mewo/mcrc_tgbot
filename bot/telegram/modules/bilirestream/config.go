package bilirestream

import "github.com/0w0mewo/mcrc_tgbot/config"

const help = `usage: 
/bilirestream gentoken: generate token, required when restreaming, the newly generated token will be replaced if it issued multiple times.
/bilirestream do <TOKEN> <BILIBILI LIVE ROOM ID>: restream	
/bilirestream status: check restream status, return "ok" if there is a restreaming task running
/bilirestream stop <TOKEN>: stop restream
------------------------------------------
Since it directly ports to the RESTful API of the backend, you can use the equvliant functions:
restream: http://mcrcsci.duckdns.org/bilirestream/restream?user=<telegram username>&token=<token generated via "/bilirestream gentoken" command>&biliroom=<room id>
check status: http://mcrcsci.duckdns.org/bilirestream/status?user=<telegram username>
stop restreaming: http://mcrcsci.duckdns.org/bilirestream/stop?user=<telegram username>&token=<token generated via "/bilirestream gentoken" command>
`

type bilirestreamConf struct {
	config config.ConfigType
}

func ConfigFrom(c config.ConfigType) *bilirestreamConf {
	cfg := &bilirestreamConf{
		config: c,
	}

	cfg.Reload()

	return cfg
}

func (mc *bilirestreamConf) Reload() {

}
