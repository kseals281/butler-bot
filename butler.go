package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"time"
)

type Butler struct {
	discord         *discordgo.Session
	activeReminders []Reminders
	villains        map[string]Villain
}

type Villain struct {
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Debut     string `json:"debut"`
	DebutDate string `json:"debut_date"`
}

type Reminders struct {
	msg  string
	time time.Time
	usr  *discordgo.User
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

//func (b *Butler) newReminder(rawReminder string) {
//
//}
