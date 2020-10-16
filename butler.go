package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	s "strings"
	"time"
)

var EmbedAuthor = &discordgo.MessageEmbedAuthor{
	Name:    "Alfred Pennyworth",
	IconURL: "https://img1.looper.com/img/uploads/2018/05/alfred-pennyworth-batman-animated-series-2.jpg",
}

type Message interface {
	createEmbed() *discordgo.MessageEmbed
}

type Butler struct {
	discord         *discordgo.Session
	activeReminders []Reminder
	villains        map[string]Villain
}

type Villain struct {
	Identity  string
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Debut     string `json:"debut"`
	DebutDate string `json:"debut_date"`
	URL       string `json:"url,omitempty"`
}

type Reminder struct {
	time    time.Time
	msg     string
	chID    string
	usr     *discordgo.User
	created time.Time
}

func (b *Butler) loadButler(s *discordgo.Session) {
	b.discord = s
	b.loadVillains()
}

func (b *Butler) loadVillains() {
	dat, err := ioutil.ReadFile("villains.json")
	panicCheck("unable to read file", err)

	err = json.Unmarshal(dat, &b.villains)
	panicCheck("error unmarshalling villain json data", err)
}

func (b *Butler) sendMessage(m Message, chID string) error {
	_, err := b.discord.ChannelMessageSendEmbed(chID, m.createEmbed())
	errCheck("error sending message", err)
	return nil
}

func (b Butler) createEmbed(msg string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:  s.Title(msg),
		Footer: &discordgo.MessageEmbedFooter{Text: "To view all commands type !commands"},
		Author: EmbedAuthor,
	}

	return embed
}

func (v Villain) createEmbed() *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       v.Name,
		Description: v.Bio,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("%s debuted in %s, %s", v.Name, v.Debut, v.DebutDate),
			IconURL: "https://cdn.shopify.com/s/files/1/1045/2900/products/Batman-Symbol_grande.png?v=1587567739",
		},
		Image: &discordgo.MessageEmbedImage{
			URL: v.URL,
		},
		Thumbnail: nil,
		Author:    EmbedAuthor,
	}

	return embed
}

func (r Reminder) createEmbed() *discordgo.MessageEmbed {
	embed := new(discordgo.MessageEmbed)
	embed.Author = EmbedAuthor

	embed.Title = fmt.Sprintf("You wanted me to remind you, Master %s.", r.usr.Mention())
	embed.Description = r.msg
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: "https://i.pinimg.com/600x315/10/aa/90/10aa90c7fab8fafd82abfeefdde4ceec.jpg",
	}

	return embed
}
