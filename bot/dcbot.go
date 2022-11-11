package bot

import (
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type DcMsgHandler func(s *discordgo.Session, m *discordgo.MessageCreate)

// TODO
type DiscordBot struct {
	logger *logrus.Entry
	bot    *discordgo.Session
}

func NewDiscordBot(token string) (*DiscordBot, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	logger := utils.NewLogger().WithField("service", "dcbot")

	return &DiscordBot{
		logger: logger,
		bot:    bot,
	}, nil
}

func (dc *DiscordBot) Start() {
	dc.bot.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentDirectMessages

	// load modules
	mods := DcModRegister.GetModules()

	for _, m := range mods {
		m.Start(dc)
	}

	// register handlers
	for _, h := range DcModRegister.GetDcCmdHandlers() {
		dc.bot.AddHandler(h)
	}

	err := dc.bot.Open()
	if err != nil {
		dc.logger.Error(err)

	}
}

func (dc *DiscordBot) Stop() {
	mods := DcModRegister.GetModules()

	for _, m := range mods {
		m.Stop(dc)
	}
	dc.bot.Close()
}

func (dc *DiscordBot) Bot() *discordgo.Session {
	return dc.bot
}
