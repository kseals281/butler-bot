package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	s "strings"

	"github.com/bwmarrin/discordgo"
)

const CommandPrefix = "!"

func (b *Butler) CommandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	botID := discord.State.User.ID
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}
	member, err := discord.GuildMember(message.GuildID, message.Author.ID)
	if err != nil {
		log.Printf("unable to get guild information for author: %+v", err)
	}

	content := message.Content

	switch {

	case s.HasPrefix(content, CommandPrefix+"hello"):
		b.hello(member, message.ChannelID)

	case s.HasPrefix(content, CommandPrefix+"commands"):
		b.commands(member, message.ChannelID)

	case s.HasPrefix(content, CommandPrefix+"oof"):
		f, err := os.Open("oof.png")
		if err != nil {
			errCheck("Something went wrong. Unable to open oof file at this time", err)
		} else {
			defer f.Close()
		}
		ms := &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "commands/pics/oof.png",
					Reader: f,
				},
			},
		}
		_, err = discord.ChannelMessageSendComplex(message.ChannelID, ms)
		if err != nil {
			errCheck("Unable to send oof to channel", err)
		}

	case s.HasPrefix(content, CommandPrefix+"villains"):
		b.knownVillains(message.ChannelID)

	case s.HasPrefix(content, CommandPrefix+"bio"):
		b.biography(s.TrimPrefix(content, CommandPrefix+"bio"), message.ChannelID)

	}
}

func (b *Butler) hello(member *discordgo.Member, chID string) {
	m := fmt.Sprintf("Good evening master %s.", member.Nick)
	_, err := b.sendMessage(String(m), chID)
	errCheck("", err)
}

func (b *Butler) commands(member *discordgo.Member, chID string) {
	commands := []*discordgo.MessageEmbedField{
		{
			Name:   "bio",
			Value:  "`bio <name>`\tLearn about a character",
			Inline: false,
		}, {
			Name:   "commands",
			Value:  "Get a list of commands",
			Inline: false,
		}, {
			Name:   "hello",
			Value:  "Receive a greeting",
			Inline: false,
		}, {
			Name:   "oof",
			Value:  fmt.Sprintf("Big oof, Master %s", member.Nick),
			Inline: false,
		}, {
			Name:   "villains",
			Value:  "List all currently known villains",
			Inline: false,
		},
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Commands",
		Description: fmt.Sprintf("All phrases I will respond to if prefixed by %s", CommandPrefix),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Please leave all suggestions in butler-suggestions",
			IconURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Question_mark_%28black%29.svg/200px-Question_mark_%28black%29.svg.png",
		},
		Author: EmbedAuthor,
		Fields: commands,
	}

	_, err := b.discord.ChannelMessageSendEmbed(chID, embed)
	errCheck("unable to send commands", err)
}

func (b *Butler) knownVillains(chID string) {
	var villainNames []string

	for name := range b.villains {
		villainNames = append(villainNames, s.TrimSpace(name))
	}
	sort.Strings(villainNames)

	var embedVillainNames []*discordgo.MessageEmbedField

	for i, n := range villainNames {
		curr := &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%d.", i+1),
			Value:  s.Title(n),
			Inline: true,
		}
		embedVillainNames = append(embedVillainNames, curr)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Villains",
		Description: "Here is a list of all currently known villains.",
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "To view all commands type !commands",
			IconURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Question_mark_%28black%29.svg/200px-Question_mark_%28black%29.svg.png",
		},
		Author: EmbedAuthor,
		Fields: embedVillainNames,
	}
	_, err := b.discord.ChannelMessageSendEmbed(chID, embed)
	errCheck("error sending list of villains", err)
}

func (b *Butler) biography(name string, chID string) {
	name = s.TrimSpace(name)
	villain, ok := b.villains[s.ToLower(name)]
	if !ok {
		return
	}
	_, err := b.sendMessage(villain, chID)
	errCheck("error sending bio", err)
}
