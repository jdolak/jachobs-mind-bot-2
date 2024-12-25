package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// relpresents an "owes" relationhip
// e.g. user1-user2 = user1 owes user2
type Db_debt struct {
	gorm.Model
	UidRelation string
	Lender      string
	Borrower    string
	Venmo       string
	Owes        int
}

func split(s *discordgo.Session, i *discordgo.InteractionCreate) {

	//values := debt_parse(s, i)
	// TODO : same as owes but subtract query ammount is:
	//query_amount, err := strconv.Atoi(values[len(values)-1])
	//checkErr(err)
	//indiv_amount := query_amount/len(values) - 1

}
func owes(s *discordgo.Session, i *discordgo.InteractionCreate) {

	response := ""

	caller := i.Member.User.ID
	values := debt_parse(s, i)
	query_amount, err := strconv.Atoi(values[len(values)-1])
	checkErr(err)

	for _, u := range values {
		relation := debt_relation(caller, u)
		inverse := debt_relation(u, caller)

		amount := debt_query(relation) - debt_query(inverse) + query_amount

		response = debt_response(amount, relation, inverse, caller, u, response)
	}
	slash_response(s, i, response)
}
func loan(s *discordgo.Session, i *discordgo.InteractionCreate) {

	response := ""

	caller := i.Member.User.ID
	values := debt_parse(s, i)
	query_amount, err := strconv.Atoi(values[len(values)-1])
	checkErr(err)

	for _, u := range values {

		relation := debt_relation(u, caller)
		inverse := debt_relation(caller, u)

		amount := debt_query(relation) - debt_query(inverse) + query_amount

		response = debt_response(amount, relation, inverse, caller, u, response)

	}

}
func paid(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO : same as owes but subtract query ammount
}
func recieved(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// TODO : same as claim but subtract query ammount
}

func register(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func links(s *discordgo.Session, i *discordgo.InteractionCreate) {

	caller := i.Member.User.ID
	var db_debt []Db_debt

	response := ""

	err := db.Where("borrower = ? AND owes != 0", caller).Find(&db_debt).Error
	checkErr(err)

	for _, u := range db_debt {
		response = response + "https://venmo.com/pay/" + u.Venmo + "?amount=" + strconv.Itoa(u.Owes) + "\n"
	}
}

func debt_parse(s *discordgo.Session, i *discordgo.InteractionCreate) []string {

	var values []string

	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Type == discordgo.ApplicationCommandOptionUser {
			values = append(values, opt.UserValue(s).ID)
		} else if opt.Type == discordgo.ApplicationCommandOptionInteger {
			values = append(values, strconv.Itoa(int(opt.IntValue())))
		} else {
			values = append(values, opt.StringValue())
		}
	}

	return values
}

func debt_query(relation string) int {
	var db_debt Db_debt

	result := db.First(&db_debt, "UidRelation = ?", relation)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		db.Create(&Db_debt{UidRelation: relation, Owes: 0})

		result = db.First(&db_debt, "uid = ?", relation)
		checkErr(result.Error)

	} else if result.Error != nil {
		checkErr(result.Error)
	} else {
		return db_debt.Owes
	}

	return db_debt.Owes
}

func debt_update(relation string, amount int) int {
	var db_debt Db_debt

	result := db.First(&db_debt, "UidRelation = ?", relation)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		db.Create(&Db_debt{UidRelation: relation, Owes: 0})

		result = db.First(&db_debt, "uid = ?", relation)
		checkErr(result.Error)

	} else if result.Error != nil {
		checkErr(result.Error)
	}

	db.Model(&db_debt).Update("owes", amount)
	return amount
}

func debt_relation(user1 string, user2 string) string {
	user1 = strings.Trim(user1, "<>@")
	user2 = strings.Trim(user2, "<>@")

	return user1 + "-" + user2
}

func debt_response(amount int, relation string, inverse string, caller string, user string, response string) string {
	if amount < 0 {
		net_amount := amount * -1
		debt_update(relation, 0)
		debt_update(inverse, net_amount)
		response = response + uwrap(caller) + " now owes " + strconv.Itoa(net_amount) + " to " + uwrap(user) + "\n"
	} else if amount == 0 {
		debt_update(relation, 0)
		debt_update(inverse, 0)
		response = response + uwrap(caller) + " and " + uwrap(user) + " are even" + "\n"
	} else {
		debt_update(relation, 0)
		debt_update(inverse, amount)
		response = response + uwrap(user) + " now owes " + strconv.Itoa(amount) + " to " + uwrap(caller) + "\n"
	}

	return response
}
