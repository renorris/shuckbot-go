package commands

import (
	"github.com/bwmarrin/discordgo"
)

func AvatarHandler(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) *HandlerResponse {
	avatarURL := ""
	if len(parameters) > 0 && len(m.Mentions) > 0 {
		avatarURL = m.Mentions[0].AvatarURL("256")
	} else {
		avatarURL = m.Author.AvatarURL("256")
	}

	response := new(HandlerResponse)
	response.AddReplyMessage(avatarURL)
	return response
}
