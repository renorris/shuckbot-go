package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token  string = os.Getenv("DISCORD_TOKEN")
	Prefix string = os.Getenv("COMMAND_PREFIX")
)

func main() {
	fmt.Println("Starting shuckbot-go")

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating session: ", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	fmt.Println("Running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("\nShutting down...")

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore non-commands
	if !strings.HasPrefix(m.Content, Prefix) {
		return
	}

	// Signal typing
	s.ChannelTyping(m.ChannelID)

	// Log message
	fmt.Println("Got command \"", m.Content, "\" from", m.Author.ID, m.Author.Username)

	// Process command
	handlerResponse := ProcessCommand(s, m)

	// Send replies
	for _, msg := range handlerResponse.GetReplyMessages() {
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			fmt.Println("There was an error sending a message: ", err)
		}
	}
}
