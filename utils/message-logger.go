package utils

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func find_channel_by_id(s *discordgo.Session, guildID string, channelName string) (string, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Name == channelName {
			return channel.ID, nil
		}
	}

	return "", fmt.Errorf("Nie znaleziono kanału: %s", channelName)
}

func message_logger(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	logChannelName := "logi-wiadomości"

	logChannelID, err := find_channel_by_id(s, m.GuildID, logChannelName)
	if err != nil {
		fmt.Println("Nie znaleziono kanału")
		return
	}

	logMessage := "Wiadomość od " + m.Author.Username + " w kanale <#" + m.ChannelID + ">: " + m.Content

	_, err = s.ChannelMessage(logChannelID, logMessage)
	if err != nil {
		log.Println("Nie udało się wysłać wiadomości")
	}
}
