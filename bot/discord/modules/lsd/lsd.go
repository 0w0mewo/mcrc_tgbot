package lsd

import (
	dcbot "github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const modname = "dc.mod.lsd"

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.GetConfigFile())
	config.RegisterModuleConfig(modname, cfg)

	// load module
	m := &lsd{
		logger: utils.NewLogger(),
		conf:   cfg,
	}
	dcbot.DcModRegister.RegistryMod(m)

}

type lsd struct {
	dcbot   *discordgo.Session
	conf    *lsdConf
	logger  *logrus.Logger
	running bool
}

func (ma *lsd) Start(b *dcbot.DiscordBot) {
	if !ma.running {
		ma.dcbot = b.Bot()
		ma.running = true
	}

	ma.Reload()

	ma.logger.Printf("%s loaded", ma.Name())
}

func (ma *lsd) Name() string {
	return modname
}

func (ma *lsd) Stop(b *dcbot.DiscordBot) {
	ma.running = false
	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *lsd) Reload() {

}
