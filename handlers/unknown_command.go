package handlers

import "github.com/bwmarrin/discordgo"

func UnknownCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Unknown command")
}
