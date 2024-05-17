package ping

import (
	"github.com/bwmarrin/discordgo"
)

// Command defines the 'ping' command
var Command = &discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "This is to check connection with the bot",
}

// Handle executes the command logic
func Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong! You have successfully established connection with the bot.",
		},
	})
}
