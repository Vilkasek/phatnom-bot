package commands

import (
	"log"

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

func Setup_Command(s *discordgo.Session, m *discordgo.MessageCreate) {
}
