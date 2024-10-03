package main

import (
	"fmt"
	"log"
	"phantom-bot/commands"
	"phantom-bot/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config, err := Load_Config("./config.json")
	if err != nil {
		log.Fatalf("Nie udało się załadować tokenu: %v", err)
	}

	bot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("Nie udało się zainicjalizować bota: %v", err)
	}

	bot.AddHandler(commands.Setup)
	bot.AddHandler(utils.Logger)

	err = bot.Open()
	if err != nil {
		log.Fatalf("Nie udało się utworzyć websocketu: %v", err)
	}

	fmt.Println(
		"Udało sie utworzyć bota oraz połączyć go.\nW celu konserwacji kodu należy bota wyłączyć, ale najpierw poinformuj użytkowników.\nAby wyłączyć bota należy użyć kombinacji klawiszy CTRL + C",
	)

	defer bot.Close()

	select {}
}
