package bot

import (
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type DcMsgHandler func(s *discordgo.Session, m *discordgo.MessageCreate)
type DiscordSlashCmdEntry struct {
	handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
	descriptor *discordgo.ApplicationCommand
}

func NewDiscordSlashCmdEntry(desc *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) *DiscordSlashCmdEntry {
	return &DiscordSlashCmdEntry{
		handler:    handler,
		descriptor: desc,
	}
}

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

	logger := utils.GetLogger().WithField("service", "dcbot")

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

	// register message handlers
	for _, h := range DcModRegister.GetDcMesgHandlers() {
		dc.bot.AddHandler(h)
	}

	// register slash command handler
	dc.bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmd := i.ApplicationCommandData().Name

		entry := DcModRegister.GetDcCmdHandler(cmd)
		if entry == nil {
			dc.logger.Error("no such command:", cmd)
			return
		}

		entry.handler(s, i)

	})

	err := dc.bot.Open()
	if err != nil {
		dc.logger.Error(err)
		return
	}

	// slash commands register on each joined group
	joinedgrps := dc.bot.State.Guilds
	appid := dc.bot.State.User.ID
	for _, grp := range joinedgrps {
		dc.logger.Infof("joined group/channel: %s-%s", grp.ID, grp.Name)

		for cmd, entry := range DcModRegister.GetDcCmdHandlers() {
			_, err = dc.bot.ApplicationCommandCreate(appid, grp.ID, entry.descriptor)
			if err != nil {
				dc.logger.Error(err)
			}
			dc.logger.Infof("register slash command: %s to guild %s", cmd, grp.Name)
		}

	}

	err = dc.bot.UpdateGameStatus(0, "f,fk")
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
