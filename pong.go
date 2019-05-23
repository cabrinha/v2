package main

import "github.com/bwmarrin/discordgo"

func pingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	// If the message is "ping" reply with "Pong!"
	if m.Content == pre+"ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == pre+"pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
