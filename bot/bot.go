package bot

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/telebot.v3"
)

// place holder for implementing bots interface
var _ Bot[telebot.Bot] = &TelegramBot{}
var _ Bot[discordgo.Session] = &DiscordBot{}

type BotType interface {
	telebot.Bot | discordgo.Session
}

type Bot[T BotType] interface {
	Start()
	Stop()
	Bot() *T
}
