package registercafe

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "register-cafe",
	Description: "This command to register cafe that you have visited before",
}

func Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modals_survey_" + i.Interaction.Member.User.ID,
			Title:    "Modals survey",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "name",
							Label:       "Cafe Name",
							Style:       discordgo.TextInputShort,
							Placeholder: "For example: Kedai Ibu Sumarti",
							Required:    true,
							MaxLength:   300,
							MinLength:   10,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "link",
							Label:       "Enter Google Maps Link Address",
							Style:       discordgo.TextInputShort,
							Placeholder: "For example: Kedai Ibu Sumarti",
							Required:    false,
							MaxLength:   40,
							MinLength:   10,
						},
					},
				},
			},
		},
	})
}

// Handle modal submission for registering a cafe
func HandleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Respond to the interaction first
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thanks for submitting your responses!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Failed to send interaction response: %v", err)
		return
	}

	// Extract data from the modal
	data := i.ModalSubmitData()

	// Safer extraction of modal input values
	if row1, ok := data.Components[0].(*discordgo.ActionsRow); ok {
		if textInput, ok := row1.Components[0].(*discordgo.TextInput); ok {
			opinion := textInput.Value
			log.Printf("Opinion: %s", opinion)
		}
	}

	if row2, ok := data.Components[1].(*discordgo.ActionsRow); ok {
		if textInput, ok := row2.Components[0].(*discordgo.TextInput); ok {
			suggestions := textInput.Value
			log.Printf("Suggestions: %s", suggestions)
		}
	}

	// TODO:
	// 1. Save data to the database
}
