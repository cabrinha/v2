package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cabrinha/v2/commands/karma"
	"github.com/cabrinha/v2/commands/ping"
	log "github.com/sirupsen/logrus"

	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func init() {
	// Setup our logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Setup our config file and read it
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s", err))
	}
}

func main() {
	// init the bot
	goBot, err := discordgo.New("Bot " + viper.GetString("token"))
	if err != nil {
		log.Warn("error creating discord session: ", err)
	}

	goBot.AddHandler(messageCreate)

	err = goBot.Open()
	if err != nil {
		log.Warn("error opening connection,", err)
		return
	}

	router := exrouter.New()
	// Ping Pong
	router.On("ping", ping.PingRoute)
	router.On("pong", ping.PongRoute)
	// Karma
	router.On("karma", karma.GetKarma)
	router.OnMatch("+-", dgrouter.NewRegexMatcher(`(\-\-|\+\+)`), karma.Handler)

	goBot.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(goBot, viper.GetString("prefix"), goBot.State.User.ID, m.Message)
	})

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	goBot.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	chanName, _ := s.Channel(m.ChannelID)
	guildName, _ := s.Guild(m.GuildID)
	messageLogger := log.WithFields(log.Fields{
		"author":  m.Author.Username,
		"server":  guildName.Name,
		"channel": chanName.Name,
	})
	messageLogger.Info(m.Content)
}
