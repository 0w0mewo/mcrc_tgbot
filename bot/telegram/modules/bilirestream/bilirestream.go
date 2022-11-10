package bilirestream

import (
	"fmt"

	tgbot "github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/service/bilirestream"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

const modname = "tg.mod.bilirestream"

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.GetConfigFile())
	config.RegisterModuleConfig(modname, cfg)

	// load module
	m := &bilirestreamer{
		logger: utils.NewLogger(),
		conf:   cfg,
		br:     bilirestream.NewClient(),
	}
	tgbot.TgModRegister.RegistryMod(m)

}

type bilirestreamer struct {
	tgbot   *telebot.Bot
	conf    *bilirestreamConf
	logger  *logrus.Logger
	br      *bilirestream.ApiClient
	running bool
}

func (ma *bilirestreamer) Start(b *tgbot.TelegramBot) {
	if !ma.running {
		ma.tgbot = b.Bot()
		ma.running = true
	}

	ma.Reload()

	tgbot.TgModRegister.AddTgEventHandler("/bilirestream", ma.bilirestream)

	ma.logger.Printf("%s loaded", ma.Name())
}

func (ma *bilirestreamer) Name() string {
	return modname
}

func (ma *bilirestreamer) Stop(b *tgbot.TelegramBot) {
	ma.running = false
	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *bilirestreamer) Reload() {

}

func (ma *bilirestreamer) bilirestream(c telebot.Context) error {
	args := c.Args()

	sender := c.Sender()
	if sender == nil {
		return c.Send("unknown error")
	}

	if len(args) < 1 {
		return c.Send(help)
	}

	switch args[0] {
	case bilirestream.API_GEN_TOKEN:
		resp, err := ma.br.Gentoken(sender.Username)
		if err != nil {
			return c.Send("ERROR: " + err.Error())
		}

		return c.Send(fmt.Sprintf(`Your username is: %s
		new token is: %s`,
			sender.Username, resp.Token))

	case bilirestream.API_BILIBILI_RESTREAM_START:
		if len(args) < 3 {
			return c.Send("not enough arguments")
		}

		resp, err := ma.br.BiliRestream(bilirestream.API_BILIBILI_RESTREAM_START, sender.Username, args[1], args[2])
		if err != nil {
			return c.Send("ERROR: " + err.Error())
		}
		if resp.Code != 0 {
			return c.Send("ERROR: "+resp.ErrMsg, ", it might be a token error")
		}

		return c.Send("Restreaming to: " + resp.Url)

	case bilirestream.API_BILIBILI_RESTREAM_STOP:
		if len(args) < 2 {
			return c.Send("not enough arguments")
		}

		resp, err := ma.br.BiliRestream(bilirestream.API_BILIBILI_RESTREAM_STOP, sender.Username, args[1], "")
		if err != nil {
			return c.Send("ERROR: " + err.Error())
		}
		if resp.Code != 0 {
			return c.Send("ERROR: " + resp.ErrMsg)
		}

		return c.Send("OK")

	case bilirestream.API_BILIBILI_RESTREAM_STATUS:
		resp, err := ma.br.BiliRestream(bilirestream.API_BILIBILI_RESTREAM_STATUS, sender.Username, "", "")
		if err != nil {
			return c.Send("ERROR: " + err.Error())
		}
		if resp.Code != 0 {
			return c.Send("ERROR: " + resp.ErrMsg)
		}
		return c.Send("RUNNING")

	default:
		return c.Send("unknown command")
	}

}
