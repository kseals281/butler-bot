package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
)

var EmbedAuthor = &discordgo.MessageEmbedAuthor{
	Name:    "Alfred Pennyworth",
	IconURL: "https://img1.looper.com/img/uploads/2018/05/alfred-pennyworth-batman-animated-series-2.jpg",
}

type String string

type Message interface {
	createEmbed() *discordgo.MessageEmbed
}

type Butler struct {
	discord  *discordgo.Session
	villains map[string]Villain
}

type Villain struct {
	Identity  string
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Debut     string `json:"debut"`
	DebutDate string `json:"debut_date"`
	URL       string `json:"url,omitempty"`
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

func (b *Butler) member(guildID, userID string) *discordgo.Member {
	member, err := b.discord.GuildMember(guildID, userID)
	errCheck("unable to get guild member data", err)
	return member
}

func (b *Butler) sendMessage(m Message, chID string) (*discordgo.Message, error) {
	embed := m.createEmbed()
	embed.Author = EmbedAuthor
	msg, err := b.discord.ChannelMessageSendEmbed(chID, embed)
	errCheck("error sending message", err)
	return msg, nil
}

func (s String) createEmbed() *discordgo.MessageEmbed {
	embed := new(discordgo.MessageEmbed)
	embed.Description = string(s)
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "Alfred is a work in progress.",
	}
	return embed
}

func (v Villain) createEmbed() *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author:      EmbedAuthor,
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
	}

	return embed
}
