package main

import (
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
)

func main() {

	dotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	appId := os.Getenv("APPID")
	guildID := os.Getenv("GUILDID")

	d, err := discordgo.New("Bot " + token)
	checkErr(err)

	d.AddHandler(messageCreate)
	d.Identify.Intents = discordgo.IntentsGuildMessages

	register_slash_commands(d, appId, guildID)

	err = d.Open()
	checkErr(err)
	set_status(d)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	d.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	} else if credit_perm_check(s, m) {
		return
	}

	autoMemes(s, m)

	if strings.HasPrefix(m.Content, "$") {
		words := strings.Split(m.Content, " ")

		switch words[0] {
		case "$credit":
			credit(s, m, nil)
		case "$debug":
			s.ChannelMessageSend(m.ChannelID, "```\n"+m.Author.ID+"\n"+m.Content+"\n```")
		}
	}

}

func autoMemes(s *discordgo.Session, m *discordgo.MessageCreate) {

	auto_words := []string{"1984"}

	for _, word := range auto_words {
		if strings.Contains(m.Content, word) {

			filepath := MEME_FOLDER + word + ".gif"
			file, err := os.Open(filepath)
			checkErr(err)

			_, err = s.ChannelFileSend(m.ChannelID, filepath, file)
			checkErr(err)

		}
	}
}

func register_slash_commands(d *discordgo.Session, appId string, guildId string) {
	_, err := d.ApplicationCommandBulkOverwrite(appId, guildId, []*discordgo.ApplicationCommand{
		{
			Name:        "rant",
			Description: "Automatically creates a publish ready argument",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "topic",
					Description: "Describe what has gotten you riled up.",
					Required:    true,
				},
			},
		},
		{
			Name:        "credit",
			Description: "Keeps things civilized.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "target",
					Description: "Who would you like praise or punish.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "credit-amount",
					Description: "be sparing...",
					Required:    false,
				},
			},
		},
	})
	checkErr(err)
	d.AddHandler(slash_handler)
}

func slash_handler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	data := i.ApplicationCommandData()

	switch data.Name {

	case "rant":
		slash_response(s, i, "hello")
	case "credit":
		credit(s, nil, i)
	}
}

func slash_response(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	},
	)
}

func set_status(s *discordgo.Session) {
	err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "online",
		Activities: []*discordgo.Activity{
			{
				Name: "you...",
				Type: discordgo.ActivityTypeWatching,
			},
		},
	})
	checkErr(err)
}
