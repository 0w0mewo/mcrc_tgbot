package commands

import (
	"github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/bot/modules/commands/cmds"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.Config.GetConfigFile())
	config.Config.RegisterModuleConfig("mod.commander", cfg)

	// load module
	cmder := &Commander{
		logger: utils.NewLogger(),
		conf:   cfg,
		cman:   newCmdManager(),
	}
	bot.ModRegister.RegistryMod(cmder)

}

// group message random repeater
type Commander struct {
	tgbot   *telebot.Bot
	conf    *commanderConf
	logger  *logrus.Logger
	cman    *cmdMan
	running bool
}

func (c *Commander) Start(b *bot.Bot) {
	if !c.running {
		c.tgbot = b.Bot()
		c.running = true
	}

	// register commands here
	c.cman.Registry("/start", cmds.Start)
	c.cman.Registry("/lsd", cmds.Lsd)
	c.cman.Registry("/help", cmds.Help)

	c.Reload()

	// load command handlers
	regtable := c.cman.GetRegTable()
	for cmd, handler := range regtable {
		bot.ModRegister.AddTgEventHandler(cmd, telebot.HandlerFunc(handler))
	}

	c.logger.Printf("%s loaded", c.Name())
}

func (c *Commander) Name() string {
	return "mod.commander"
}

func (c *Commander) Stop(b *bot.Bot) {
	c.running = false
	c.logger.Printf("%s unloaded", c.Name())

}

func (c *Commander) Reload() {

}
