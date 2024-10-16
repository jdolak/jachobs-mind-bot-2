package main

import (
	"errors"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bwmarrin/discordgo"
)

type Db_credit struct {
	gorm.Model
	Uid    string
	Credit int
}

var db *gorm.DB = create_db()

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

func credit(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate) {
	var words []string
	response := ""

	if m != nil {
		words = strings.Split(m.Content, " ")
	} else {
		words = append(words, "$credit")
		options := i.ApplicationCommandData().Options
		for _, opt := range options {
			if opt.Type == discordgo.ApplicationCommandOptionUser {
				words = append(words, opt.UserValue(s).ID)
			} else if opt.Type == discordgo.ApplicationCommandOptionInteger {
				words = append(words, strconv.Itoa(int(opt.IntValue())))
			} else {
				words = append(words, opt.StringValue())
			}
		}
	}

	if len(words) == 3 {
		var author string
		if m != nil {
			author = m.Author.ID
		} else {
			author = i.Member.User.ID
		}
		user_data := credit_check_user(words[1])
		caller := credit_check_user(author)
		num, err := strconv.Atoi(words[2])
		checkErr(err)

		if i != nil || !check_role(m, JACHOB_ROLE) {
			if user_data.Uid == author {
				response = response + "no. -50 credit.\n"
				num = -50
			} else if num >= 500 || num <= -500 {
				response = response + "no. -50 credit for u.\n"
				user_data = caller
				num = -50
			} else if caller.Credit <= 0 {
				response = response + "only people in good standing can give credit. -1 credit."
				user_data = caller
				num = -1
			}
		}

		credit := user_data.Credit + num
		db.Model(&user_data).Update("credit", credit)
		user_data = credit_check_user(user_data.Uid)
		response = response + uwrap(user_data.Uid) + " now has " + strconv.Itoa(user_data.Credit) + " credit"

		if m != nil {
			s.ChannelMessageSend(m.ChannelID, response)
		} else {
			slash_response(s, i, response)
		}

	} else if len(words) == 2 {
		user_data := credit_check_user(words[1])
		response := uwrap(user_data.Uid) + " has " + strconv.Itoa(user_data.Credit) + " credit"
		if m != nil {
			s.ChannelMessageSend(m.ChannelID, response)
		} else {
			slash_response(s, i, response)
		}

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

func credit_perm_check(s *discordgo.Session, m *discordgo.MessageCreate) bool {

	user := credit_check_user(m.Author.ID)

	if !check_role(m, JACHOB_ROLE) {
		if user.Credit <= 0 && strings.Contains(m.Content, "http") {
			s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
			s.ChannelMessageSend(m.ChannelID, "uhhhh... "+uwrap(m.Author.ID)+" check ur credit.")

			return true
		}
	}
	return false
}

func create_db() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("/botdata/data.db"), &gorm.Config{})
	checkErr(err)

	db.AutoMigrate(&Db_credit{})

	return db
}
