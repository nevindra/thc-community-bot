package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	cafe "github.com/nevindra/community-bot/cafe"
	"github.com/nevindra/community-bot/command"
	setup "github.com/nevindra/community-bot/setup"
)

var (
	BotToken string
	s        *discordgo.Session
)

type ModalHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

func main() {
	db, err := setup.DB()
	if err != nil {
		log.Fatalf("cannot connect to the database: %v", err)
	}
	defer db.Close()

	cafeDomain, err := cafe.NewCafeDomain(db)
	if err != nil {
		log.Fatalf("cannot create cafe domain: %v", err)
	}
	commandDomain, err := command.NewCommandDomain(cafeDomain)
	if err != nil {
		log.Fatalf("cannot create command domain: %v", err)
	}

	BotToken := os.Getenv("BOT_TOKEN")
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	// Make sure to add this handler each time you add new command
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		command.Ping.Name:         commandDomain.HandlePing,
		command.Event.Name:        commandDomain.HandleEvent,
		command.RegisterCafe.Name: commandDomain.HandleRegisterCafe,
		command.ListCafe.Name:     commandDomain.HandleListCafe,
		// Add other commands as necessary
	}

	// Map to hold modal handlers, Make sure to add this handler each time you add new command
	modalHandlers := map[string]ModalHandler{
		"register_cafe_": commandDomain.HandleRegisterCafeModalSubmit,
		"create_event":   commandDomain.HandleEventModalSubmit,
		"list_cafe_":     commandDomain.HandleListCafeModalSubmit,
		// Add other modal handlers as necessary
	}

	// Register command handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionModalSubmit:
			modalID := i.ModalSubmitData().CustomID
			// Find the correct handler based on a prefix or pattern
			found := false
			for prefix, handler := range modalHandlers {
				if strings.HasPrefix(modalID, prefix) {
					handler(s, i)
					found = true
					break
				}
			}
			if !found {
				log.Printf("No handler found for modal ID: %s", modalID)
			}
		}
	})

	// Log bot login details
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open the Discord session
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	// Register commands with Discord
	commands := []*discordgo.ApplicationCommand{
		command.Ping,
		command.RegisterCafe,
		command.Event,
		command.ListCafe,
		// Add other commands
	}

	guildID := "" // If you want to register commands globally, set guildID to an empty string
	registeredCommands, err := registerCommands(s, commands, guildID)
	if err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}
	defer unregisterCommands(s, registeredCommands, guildID)

	// Wait for a signal to exit
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Bot is now running. Press Ctrl+C to exit.")
	<-stop

	log.Println("Bot is shutting down...")
}

// registerCommands registers the given list of commands for the application
func registerCommands(s *discordgo.Session, commands []*discordgo.ApplicationCommand, guildID string) ([]*discordgo.ApplicationCommand, error) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, cmd := range commands {
		registeredCmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
		if err != nil {
			return nil, fmt.Errorf("cannot create '%v' command: %v", cmd.Name, err)
		}
		registeredCommands[i] = registeredCmd
	}
	return registeredCommands, nil
}

// unregisterCommands unregisters commands on bot shutdown
func unregisterCommands(s *discordgo.Session, commands []*discordgo.ApplicationCommand, guildID string) {
	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
		if err != nil {
			log.Printf("Failed to delete '%v' command: %v", cmd.Name, err)
		}
	}
}
