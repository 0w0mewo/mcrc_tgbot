package repeater

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
	config.Config.RegisterModuleConfig("mod.repeater", cfg)

	// load module
	r := &Repeater{
		logger:        utils.NewLogger(),
		conf:          cfg,
		chatRandLimit: utils.NewRandomMap(cfg.randstart, cfg.randend),
		msgCounter:    utils.NewCounter(),
	}
	bot.ModRegister.Registry(r)
}

// group message random repeater
type Repeater struct {
	tgbot         *telebot.Bot
	conf          *repeaterConf
	chatRandLimit *utils.RandomMap
	msgCounter    *utils.Counter
	logger        *logrus.Logger
	running       bool
}

func (r *Repeater) Start(b *bot.Bot) {
	if !r.running {
		r.tgbot = b.Bot()
		r.running = true
	}

	r.Reload()
	r.logger.Printf("%s loaded", r.Name())
}

func (r *Repeater) Name() string {
	return "mod.repeater"
}

func (r *Repeater) Stop(b *bot.Bot) {
	r.running = false
	r.logger.Printf("%s unloaded", r.Name())
}

func (r *Repeater) Reload() {
	r.chatRandLimit.SetStart(r.conf.randstart)
	r.chatRandLimit.SetEnd(r.conf.randend)

	// set repeater on text and sticker message
	for _, ev := range []string{telebot.OnText, telebot.OnSticker} {
		if r.tgbot == nil {
			return
		}
		r.tgbot.Handle(ev, r.repeater, bot.MessageCounter(r.msgCounter, r.chatRandLimit))
	}

}

func (r *Repeater) repeater(c telebot.Context) error {
	if !c.Message().FromGroup() {
		return nil
	}

	// text message
	if txtMsg := c.Message().Text; txtMsg != "" {
		return c.Send(txtMsg)
	}

	// sticker message
	if stickerMsg := c.Message().Sticker; stickerMsg != nil {
		return c.Send(stickerMsg)
	}

	println("called")

	// random qutoe for other type of messages
	m := utils.RandChoice(r.conf.qutoes)
	return c.Send(m)

}
