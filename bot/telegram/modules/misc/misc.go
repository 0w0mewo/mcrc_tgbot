package misc

import (
	tgbot "github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

const modname = "tg.mod.misc"

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.GetConfigFile())
	config.RegisterModuleConfig(modname, cfg)

	// load module
	m := &miscMod{
		logger: utils.GetLogger().WithField("module", modname),
		conf:   cfg,
	}
	tgbot.TgModRegister.RegistryMod(m)

}

type miscMod struct {
	tgbot   *telebot.Bot
	conf    *miscModConf
	logger  *logrus.Entry
	running bool
}

func (ma *miscMod) Start(b *tgbot.TelegramBot) {
	if !ma.running {
		ma.tgbot = b.Bot()
		ma.running = true
	}

	ma.Reload()

	tgbot.TgModRegister.AddTgEventHandler("/start", func(c telebot.Context) error {
		return c.Send("budong")
	})

	tgbot.TgModRegister.AddTgEventHandler("/help", func(c telebot.Context) error {
		return c.Send("nope")
	})

	ma.logger.Printf("%s loaded", ma.Name())
}

func (ma *miscMod) Name() string {
	return modname
}

func (ma *miscMod) Stop(b *tgbot.TelegramBot) {
	ma.running = false
	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *miscMod) Reload() {

}
