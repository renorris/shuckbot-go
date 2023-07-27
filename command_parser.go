package main

import (
	"shuckbot-go/handlers"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Parse command and return a handler function
func ParseCommand(s *discordgo.Session, m *discordgo.MessageCreate) func() {
	command := strings.TrimPrefix(m.Content, Prefix)
	//parameters := strings.Split(m.Content, " ")[1:]

	switch command {
	case "ping":
		return func() { handlers.PingHandler(s, m) }
	default:
		return func() { handlers.UnknownCommandHandler(s, m) }
	}
}
