package command

import (
	"log"

	"github.com/bwmarrin/discordgo"
	cafe "github.com/nevindra/community-bot/cafe"
)

var RegisterCafe = &discordgo.ApplicationCommand{
	Name:        "register-cafe",
	Description: "This command to register cafe that you have visited before",
}

func (c *CommandDomain) HandleRegisterCafe(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "register_cafe_" + i.Interaction.Member.User.ID,
			Title:    "Register Cafe",
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
							MinLength:   5,
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
							Required:    true,
							MaxLength:   300,
							MinLength:   10,
						},
					},
				},
			},
		},
	})
}

// Handle modal submission for registering a cafe
func (c *CommandDomain) HandleRegisterCafeModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	var (
		name, gmapslink string
	)
	// Safer extraction of modal input values
	if row1, ok := data.Components[0].(*discordgo.ActionsRow); ok {
		if textInput, ok := row1.Components[0].(*discordgo.TextInput); ok {
			name = textInput.Value
		}
	}

	if row2, ok := data.Components[1].(*discordgo.ActionsRow); ok {
		if textInput, ok := row2.Components[0].(*discordgo.TextInput); ok {
			gmapslink = textInput.Value
		}
	}

	cafeData := cafe.Cafe{
		Name:                   name,
		GmapsLink:              gmapslink,
		CreatedByDiscordUserID: i.Interaction.Member.User.ID,
	}
	err = c.cafeDomain.Insert(&cafeData)
	if err != nil {
		log.Printf("Failed to insert cafe data: %v", err)

		// Respond to the interaction first
		errMsg := "failed to insert cafe data. please contact admin"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	log.Printf("Cafe data inserted successfully")
}
