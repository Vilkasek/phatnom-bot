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

	// ID roli, która będzie miała dostęp do kanałów (zamień na prawdziwe ID roli)
	roleID := "ROLE_ID"

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
	var categoryID string
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory &&
			strings.ToLower(channel.Name) == strings.ToLower(categoryName) {
			categoryID = channel.ID
			break
		}
	}

	// Tworzenie kategorii, jeśli nie istnieje
	if categoryID == "" {
		category, err := s.GuildChannelCreate(
			m.GuildID,
			categoryName,
			discordgo.ChannelTypeGuildCategory,
		)
		if err != nil {
			log.Printf("Nie udało się stworzyć kategorii %s: %v", categoryName, err)
			return
		}
		categoryID = category.ID
		log.Printf("Utworzono kategorię %s", categoryName)
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
			"Wszystkie kanały już istnieją, pomijanie tworzenia nowych.",
		)
		return
	}

	// Tworzenie brakujących kanałów i przypisywanie do kategorii
	for _, name := range channelNames {
		lowerName := strings.ToLower(name)
		if existingChannels[lowerName] {
			log.Printf("Kanał %s już istnieje, pomijanie tworzenia...", name)
			continue // Kanał już istnieje, pomijamy jego tworzenie
		}

		// Tworzenie kanału w kategorii
		channel, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name:     name,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: categoryID, // Przypisywanie kanału do kategorii
		})
		if err != nil {
			log.Printf("Nie udało się stworzyć kanału %s: %v", name, err)
			continue
		}

		// Ustawienie uprawnień, aby tylko określona rola i administratorzy mieli dostęp do kanałów
		err = s.ChannelPermissionSet(
			channel.ID,
			roleID,
			discordgo.PermissionOverwriteTypeRole,
			discordgo.PermissionViewChannel|discordgo.PermissionSendMessages,
			0,
		)
		if err != nil {
			log.Printf("Nie udało się ustawić uprawnień dla kanału %s: %v", name, err)
			continue
		}

		// Ustawienie uprawnień dla administratorów
		err = s.ChannelPermissionSet(
			channel.ID,
			m.Author.ID,
			discordgo.PermissionOverwriteTypeMember,
			discordgo.PermissionViewChannel|discordgo.PermissionSendMessages|discordgo.PermissionAdministrator,
			0,
		)
		if err != nil {
			log.Printf(
				"Nie udało się ustawić uprawnień dla administratorów w kanale %s: %v",
				name,
				err,
			)
			continue
		}
	}

	s.ChannelMessageSend(m.ChannelID, "Kanały zostały utworzone pomyślnie w kategorii!")
}

func Setup(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!setup") {
		setup_command(s, m)
	}
}
