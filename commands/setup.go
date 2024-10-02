package commands

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func check_admin_perm(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	guildMember, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		log.Fatalf("Nie udało się pobrać informacji o użytkowniku: %v", err)
		return false
	}

	for _, roleID := range guildMember.Roles {
		role, err := s.State.Role(m.GuildID, roleID)
		if err == nil && (role.Permissions&discordgo.PermissionAdministrator != 0) {
			return true
		}
	}

	return false
}

func setup_command(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !check_admin_perm(s, m) {
		s.ChannelMessageSend(
			m.ChannelID,
			"Musisz mieć uprawnienia administratora, aby uruchomić tę komendę.",
		)
		return
	}

	// Nazwy kanałów, które mają być utworzone
	channelNames := []string{"log-kanały", "log-wiadomości"}

	// Nazwa kategorii
	categoryName := "Logi"

	// Pobranie istniejących kanałów na serwerze
	channels, err := s.GuildChannels(m.GuildID)
	if err != nil {
		log.Printf("Nie udało się pobrać kanałów serwera: %v", err)
		return
	}

	// Sprawdzanie, czy kategoria już istnieje
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory &&
			strings.ToLower(channel.Name) == strings.ToLower(categoryName) {
			break
		}
	}

	// Sprawdzanie, czy wszystkie kanały już istnieją
	allChannelsExist := true
	existingChannels := make(map[string]bool)

	for _, channel := range channels {
		existingChannels[strings.ToLower(channel.Name)] = true
	}

	for _, name := range channelNames {
		if !existingChannels[strings.ToLower(name)] {
			allChannelsExist = false
			break
		}
	}

	if allChannelsExist {
		s.ChannelMessageSend(
			m.ChannelID,
			"Wszystkie kanały już istnieją.",
		)
		return
	} else {
		s.ChannelMessageSend(m.ChannelID, "Brakuje ci kilku kanałów")
	}
}

func Setup(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!setup") {
		setup_command(s, m)
	}
}
