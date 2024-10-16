package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
)

const MEME_FOLDER = "./static/"
const JACHOB_ROLE = "1295953082827407392"

type Db_credit struct {
	gorm.Model
	Uid    string
	Credit int
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var db *gorm.DB = create_db()

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	d.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	autoMemes(s, m)

	if strings.HasPrefix(m.Content, "$") {
		words := strings.Split(m.Content, " ")

		switch words[0] {
		case "$credit":
			credit(s, m)
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
	}
}

func slash_response(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	},
	)
}

func create_db() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	checkErr(err)

	db.AutoMigrate(&Db_credit{})

	return db
}

func credit_check_user(uid string) Db_credit {
	var db_credit Db_credit

	uid_strip := strings.Trim(uid, "<>@")
	result := db.First(&db_credit, "uid = ?", uid_strip)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		db.Create(&Db_credit{Uid: uid_strip, Credit: 100})

		result = db.First(&db_credit, "uid = ?", uid_strip)
		checkErr(result.Error)

	} else if result.Error != nil {
		checkErr(result.Error)
	} else {
		return db_credit
	}

	return db_credit

}

func credit(s *discordgo.Session, m *discordgo.MessageCreate) {

	words := strings.Split(m.Content, " ")

	if len(words) == 3 {

		user_data := credit_check_user(words[1])
		caller := credit_check_user(m.Author.ID)
		num, err := strconv.Atoi(words[2])
		checkErr(err)

		if !check_role(m, JACHOB_ROLE) {
			if user_data.Uid == m.Author.ID {
				s.ChannelMessageSend(m.ChannelID, "no. -50 credit.")
				num = -50
			} else if num >= 500 || num <= -500 {
				s.ChannelMessageSend(m.ChannelID, "no. -50 credit.")
				user_data = caller
				num = -50
			} else if caller.Credit <= 0 {
				s.ChannelMessageSend(m.ChannelID, "only people in good standing can give credit. -1 credit.")
				user_data = caller
				num = -1
			}
		}

		credit := user_data.Credit + num
		db.Model(&user_data).Update("credit", credit)
		user_data = credit_check_user(words[1])

		s.ChannelMessageSend(m.ChannelID, uwrap(user_data.Uid)+" now has "+strconv.Itoa(user_data.Credit)+" credit")

	} else if len(words) == 2 {
		user_data := credit_check_user(words[1])
		s.ChannelMessageSend(m.ChannelID, uwrap(user_data.Uid)+" has "+strconv.Itoa(user_data.Credit)+" credit")

	} else {
		s.ChannelMessageSend(m.ChannelID, "Check ur syntax")
	}
}

func uwrap(uid string) string {
	return "<@" + uid + ">"
}

func check_role(m *discordgo.MessageCreate, role string) bool {
	for _, r := range m.Member.Roles {
		if r == role {
			return true
		}
	}
	return false
}
