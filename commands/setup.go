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
		s.ChannelMessageSend(m.ChannelID, "Nie masz uprawnień do wykonania tej komendy")
		return
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
