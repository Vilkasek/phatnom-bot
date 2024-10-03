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
	// Ignorujemy wiadomości wysłane przez samego bota
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Nazwa kanału, do którego mają trafiać logi wiadomości (zmień na prawdziwą nazwę kanału)
	logChannelName := "log-wiadomości"

	// Znajdowanie ID kanału na podstawie jego nazwy
	logChannelID, err := find_channel_by_id(s, m.GuildID, logChannelName)
	if err != nil {
		log.Printf("Nie udało się znaleźć kanału %s: %v", logChannelName, err)
		return
	}

	// Tworzymy wiadomość logującą, zawierającą autora, ID kanału i treść wiadomości
	logMessage := "Wiadomość od <@" + m.Author.ID + ">" + " w kanale <#" + m.ChannelID + ">: " + m.Content

	// Wysyłanie wiadomości na kanał logów
	_, err = s.ChannelMessageSend(logChannelID, logMessage)
	if err != nil {
		log.Printf("Nie udało się wysłać wiadomości do kanału logów: %v", err)
	}
}
