package command

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/nevindra/community-bot/cafe"
)

var ListCafe = &discordgo.ApplicationCommand{
	Name:        "list-cafe",
	Description: "This command to see list of cafe registered by the users here",
}

func (c *CommandDomain) HandleListCafe(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "list_cafe_" + i.Interaction.Member.User.ID,
			Title:    "List Cafe",
			Content:  "Here are the list of cafes registered by the users",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "page",
							Label:       "Page Number",
							Style:       discordgo.TextInputShort,
							Placeholder: "1",
							Required:    true,
							MaxLength:   3,
							MinLength:   1,
						},
					},
				},
			},
		},
	})
}

// Handle modal submission for listing cafes
func (c *CommandDomain) HandleListCafeModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Respond to the interaction first
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thanks for submitting your responses!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to submit your responses, please try again later",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	// Extract data from the modal
	data := i.ModalSubmitData()
	var page string
	// Safer extraction of modal input values
	if row1, ok := data.Components[0].(*discordgo.ActionsRow); ok {
		if textInput, ok := row1.Components[0].(*discordgo.TextInput); ok {
			page = textInput.Value
		}
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Printf("Failed to convert page number to integer: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid page number, please try again",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	size := 10
	offset := (pageInt - 1) * size
	cafes, err := c.cafeDomain.List(cafe.ListParam{
		Limit:  size,
		Offset: offset,
	})
	if err != nil {
		log.Printf("Failed to list cafes: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to submit your responses, please try again later",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var content string
	for i, cafe := range cafes {
		content += fmt.Sprintf("%d. %s - %s", i+1, cafe.Name, cafe.GmapsLink) + "\n"
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		log.Printf("Failed to send interaction response: %v", err)
		return
	}

	log.Print("Successfully listed cafes")
}
