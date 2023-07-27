package handlers

import "github.com/bwmarrin/discordgo"

func UnknownCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) *HandlerResponse {
	response := new(HandlerResponse)
	response.AddReplyMessage("Unknown command")
	return response
}
