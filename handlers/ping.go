package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func PingHandler(s *discordgo.Session, m *discordgo.MessageCreate) *HandlerResponse {
	response := new(HandlerResponse)
	response.AddReplyMessage("Pong!")
	return response
}
