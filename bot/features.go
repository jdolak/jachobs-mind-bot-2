package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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

func ratio(s *discordgo.Session, m *discordgo.MessageCreate) {

	message := strings.ToLower(m.Content)

	if strings.HasPrefix(message, "ratio") {
		last_msgs, err := s.ChannelMessages(m.ChannelID, 2, "", "", "")
		checkErr(err)

		if len(last_msgs) < 2 {
			return
		}

		last_msg := last_msgs[1]
		user := credit_check_user(last_msg.Author.ID)

		response := "ummm " + uwrap(user.Uid) + ", you only have " + strconv.Itoa(user.Credit) + " credit..."
		s.ChannelMessageSend(m.ChannelID, response)
		infolog.Print(m.Author.ID + " used ratio on " + user.Uid)

	}
}
