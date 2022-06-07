package misc

import (
	"github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.Config.GetConfigFile())
	config.Config.RegisterModuleConfig("mod.misc", cfg)

	// load module
	m := &miscMod{
		logger: utils.NewLogger(),
		conf:   cfg,
	}
	bot.ModRegister.RegistryMod(m)

}

type miscMod struct {
	tgbot   *telebot.Bot
	conf    *miscModConf
	logger  *logrus.Logger
	running bool
}

func (ma *miscMod) Start(b *bot.Bot) {
	if !ma.running {
		ma.tgbot = b.Bot()
		ma.running = true
	}

	ma.Reload()

	bot.ModRegister.AddTgEventHandler("/start", func(c telebot.Context) error {
		return c.Send("budong")
	})

	bot.ModRegister.AddTgEventHandler("/help", func(c telebot.Context) error {
		return c.Send("nope")
	})

	ma.logger.Printf("%s loaded", ma.Name())
}

func (ma *miscMod) Name() string {
	return "mod.misc"
}

func (ma *miscMod) Stop(b *bot.Bot) {
	ma.running = false
	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *miscMod) Reload() {

}
