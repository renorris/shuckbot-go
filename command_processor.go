package main

import (
	"strings"

	"github.com/renorris/shuckbot-go/handlers"

	"github.com/bwmarrin/discordgo"
)

// Delegate a handler & return its response
func ProcessCommand(s *discordgo.Session, m *discordgo.MessageCreate) *handlers.HandlerResponse {
	command := strings.TrimPrefix(m.Content, Prefix)
	//parameters := strings.Split(m.Content, " ")[1:]

	switch command {
	case "ping":
		return handlers.PingHandler(s, m)
	default:
		return handlers.UnknownCommandHandler(s, m)
	}
}
