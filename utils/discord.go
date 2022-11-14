package utils

import (
	"github.com/bwmarrin/discordgo"
)

func DiscordSlashCmdRespString(s *discordgo.Session, i *discordgo.Interaction, str string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: str,
		},
	})

}
