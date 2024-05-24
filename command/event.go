package command

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Command defines the 'create event' command
var Event = &discordgo.ApplicationCommand{
	Name:        "create-event",
	Description: "Create event",
}

// Handle executes the command logic
func (c *CommandDomain) HandleEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "create_event" + i.Interaction.Member.User.ID,
			Title:    "Create event",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "name",
							Label:       "Create event name",
							Style:       discordgo.TextInputShort,
							Placeholder: "For example: WFC - Jaksel",
							Required:    true,
							MaxLength:   300,
							MinLength:   5,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "descriptions",
							Label:       "Input event descriptions",
							Style:       discordgo.TextInputShort,
							Placeholder: "For example: networking",
							Required:    false,
							MaxLength:   300,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "location",
							Label:       "Send google maps link here",
							Style:       discordgo.TextInputShort,
							Placeholder: "For example: https://maps.app.goo.gl/T6vhVnvcKwYGoFag6",
							Required:    true,
							MaxLength:   30,
							MinLength:   5,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "time",
							Label:       "Set time here",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Input with DD/MM/YYYY 24 hour format. For example: 22/10/2024 18.00",
							Required:    true,
							MaxLength:   30,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to send interaction response: %v", err)
	}
}

// Handle modal submission for registering a cafe
func (c *CommandDomain) HandleEventModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var (
		name, descriptions, location, event_time_str string
	)

	// Create object that init channel name mapping the channel id to role id
	channelNameMap := make(map[string]string)
	channelNameMap["1238300702191779951"] = "1238300186384662629" // Jabodetabek
	channelNameMap["1238300702191779951"] = "1238300222002696243" // Surabaya
	channelNameMap["1159497926268162048"] = "1238300240948236318" // Malang
	channelNameMap["1159497926268162048"] = "1238300252763852841" // Bali
	channelNameMap["1159497926268162048"] = "1238300275194859601" // Bandung

	// Respond to the interaction first
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thanks for your interest in our community events. We will create your event soon.",
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
	for i, row := range data.Components {
		if actionRow, ok := row.(*discordgo.ActionsRow); ok {
			if textInput, ok := actionRow.Components[0].(*discordgo.TextInput); ok {
				if i == 0 {
					name = textInput.Value
				} else if i == 1 {
					descriptions = textInput.Value
				} else if i == 2 {
					location = textInput.Value
				} else if i == 3 {
					event_time_str = textInput.Value
				}
			}
		}
	}

	parsedTime, err := parseDateTime(event_time_str)
	if err != nil {
		errorMessage := err.Error()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errorMessage,
		})
		return
	}
	endTime := parsedTime.Add(time.Hour + 8)

	// create event based on input data
	scheduledEvent, err := s.GuildScheduledEventCreate(i.GuildID, &discordgo.GuildScheduledEventParams{
		Name:               name,
		Description:        descriptions,
		ScheduledStartTime: &parsedTime,
		ScheduledEndTime:   &endTime,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: location,
		},
	})
	if err != nil {
		log.Printf("Error creating scheduled event: %v", err)
	}

	log.Println("Created scheduled event:", scheduledEvent.Name)

	// create embed message to send in the #info-wfc
	embed := &discordgo.MessageEmbed{
		Title: "New Event Announcement!",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.User.Username,                                                                                        // Correct usage to get the username of the member
			IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=1024", i.Member.User.ID, i.Member.User.Avatar), // Correct avatar URL
		},
		Description: descriptions,
		Color:       0xD2691E, // Optional: Set embed color (replace with desired color code)
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "", // Optional: Set thumbnail image URL (replace with image URL)
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  name,
				Inline: false, // Set to true for inline fields
			},
			{
				Name:   "When",
				Value:  fmt.Sprintf("%s - %s", parsedTime.Format("02/01/2006 15.04 PM"), endTime.Format("15.04 PM")),
				Inline: false,
			},
			{
				Name:   "Location",
				Value:  location,
				Inline: false,
			},
			{
				Name:   "Join Here",
				Value:  fmt.Sprintf("https://discord.com/events/" + i.GuildID + "/" + scheduledEvent.ID),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339Nano), // Optional: Set embed timestamp
	}

	// Send the embed message to the channel based on the channel ID
	role_id := channelNameMap[i.ChannelID]

	content := fmt.Sprintf("<@&%s> Hi everyone, please check this new event!", role_id)

	_, err = s.ChannelMessageSendComplex("1240875544422125660", &discordgo.MessageSend{
		Content: content,
		Embeds:  []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		// Handle error
		log.Println("Error sending announcement:", err)
		errorMessage := err.Error()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errorMessage,
		})
		return
	}

	log.Println("Announcement sent successfully!")

	// Edit Interaction Response
	cnt := "Event created successfully! Please check the #info-wfc channel for more information."
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &cnt,
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
		return
	}
}

// parseDateTime takes a combined date and time string in "dd/mm/yyyy HH.MM" format,
// then converts the datetime to GMT+7
func parseDateTime(dateTimeStr string) (time.Time, error) {
	// Define the layout for the combined date and time format
	const layout = "02/01/2006 15.04" // dd/mm/yyyy HH.MM

	// Specify the GMT+7 location
	loc, err := time.LoadLocation("Asia/Jakarta") // Bangkok is in GMT+7
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load GMT+7 location: %w", err)
	}

	// Parse the datetime string in the specified location
	dateTime, err := time.ParseInLocation(layout, dateTimeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("please use the correct format: DD/MM/YYYY HH:MM")
	}

	// Check if the date and time are in the past
	if dateTime.Before(time.Now().In(loc)) {
		return time.Time{}, fmt.Errorf("cannot schedule event in the past, please check your input here: %s", dateTimeStr)
	}

	return dateTime, nil
}
