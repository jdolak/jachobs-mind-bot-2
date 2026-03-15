package main

import "github.com/bwmarrin/discordgo"

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
		{
			Name:        "leaderboard",
			Description: "Who's at the top...and bottom",
		},
		{
			Name:        "owes",
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
					Name:        "debt-amount",
					Description: "be sparing...",
					Required:    false,
				},
			},
		},
		{
			Name:        "loan",
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
		{
			Name:        "links",
			Description: "Who's at the top...and bottom",
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
	case "leaderboard":
		leaderboard(s, i)
	case "owes":
		owes(s, i)
	case "loan":
		loan(s, i)
	case "links":
		links(s, i)
	}
}

func slash_response(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	},
	)
}
