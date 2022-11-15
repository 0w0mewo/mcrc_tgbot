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

	dc.bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		me := r.User
		for _, grp := range r.Guilds {
			dc.logger.Infof("%s(%s) is ready to serve on channel: %s(%s)", me.Username, me.ID, grp.Name, grp.ID)
		}

		dc.logger.Infof("you can now join the bot to a server via https://discord.com/oauth2/authorize?client_id=%s&permissions=0&scope=bot", me)

	})

	// bot start
	err := dc.bot.Open()
	if err != nil {
		dc.logger.Error(err)
		return
	}

	// slash commands register
	appid := dc.bot.State.User.ID
	for cmd, entry := range DcModRegister.GetDcCmdHandlers() {
		_, err = dc.bot.ApplicationCommandCreate(appid, "", entry.descriptor)
		if err != nil {
			dc.logger.Error(err)
		}
		dc.logger.Infof("register slash command: %s", cmd)
	}

	err = dc.bot.UpdateGameStatus(0, "f,fk")
	if err != nil {
		dc.logger.Error(err)
	}
}

func (dc *DiscordBot) Stop() {
	defer func() {
		dc.bot.Close()
		dc.logger.Infof("stopped")
	}()

	mods := DcModRegister.GetModules()
	for _, m := range mods {
		m.Stop(dc)
	}

	// delete all slash commands
	appid := dc.bot.State.User.ID
	apps, err := dc.bot.ApplicationCommands(appid, "")
	if err != nil {
		dc.logger.Error(err)
		return
	}
	for _, slashcmd := range apps {
		err := dc.bot.ApplicationCommandDelete(appid, slashcmd.GuildID, slashcmd.ID)
		if err != nil {
			dc.logger.Error(err)
		}

	}

}

func (dc *DiscordBot) Bot() *discordgo.Session {
	return dc.bot
}
