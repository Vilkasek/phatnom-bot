package utils

import "github.com/bwmarrin/discordgo"

func Logger(s *discordgo.Session, m *discordgo.MessageCreate) {
	message_logger(s, m)
}
