package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	s "strings"

	"github.com/bwmarrin/discordgo"
)

const custom = "1/2/06 3:04pm"

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

	commandPrefix := "!"
	content := message.Content

	switch {

	case s.HasPrefix(content, commandPrefix+"hello"):
		b.hello(member, message.ChannelID)

	case s.HasPrefix(content, commandPrefix+"commands"):
		_, err := discord.ChannelMessageSend(message.ChannelID,
			"__**Command List**__\n"+
				"`hello: Alfred returns a greeting`\n"+
				"`commands: Alfred returns all valid commands\n"+
				"`oof: Alfred replies with a big oof`\n"+
				"`remindMe: Butler-Bot takes in a reminder for a set date and time. Format must strictly follow this example: *your message here* - Jan 2, 2006 at 3:04pm (MST)`\n")
		errCheck("Failed to send list of commands", err)

	case s.HasPrefix(content, commandPrefix+"oof"):
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

	case s.HasPrefix(content, commandPrefix+"villains"):
		b.knownVillains(message.ChannelID)

	case s.HasPrefix(content, commandPrefix+"bio"):
		b.biography(s.TrimPrefix(content, commandPrefix+"bio"), message.ChannelID)
	}
}

func (b *Butler) hello(member *discordgo.Member, chID string) {
	m := fmt.Sprintf("Good evening master %s.", member.Nick)
	_, err := b.discord.ChannelMessageSend(chID, m)
	errCheck("", err)
}

func (b *Butler) knownVillains(chID string) {
	var villainNames []string

	for name := range b.villains {
		villainNames = append(villainNames, s.TrimSpace(name))
	}
	sort.Strings(villainNames)

	fmtVillainNames := ""

	for i, n := range villainNames {
		if i == 0 {
			fmtVillainNames = s.Title(n)
			continue
		}
		fmtVillainNames += fmt.Sprintf(", %s", s.Title(n))
	}

	msg := fmt.Sprintf("Here is a list of all currently known villians: %+v", fmtVillainNames)
	_, err := b.discord.ChannelMessageSend(chID, msg)
	errCheck("error sending list of villains", err)
}

func (b *Butler) biography(name string, chID string) {
	name = s.TrimSpace(name)
	villain, ok := b.villains[s.ToLower(name)]
	if !ok {
		return
	}
	err := b.sendMessage(villain, chID)
	errCheck("error sending bio", err)
}
