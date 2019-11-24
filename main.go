package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/cabrinha/v2/commands/karma"
	"github.com/cabrinha/v2/commands/ping"
	"github.com/cabrinha/v2/commands/quotes"
	"github.com/cabrinha/v2/plugins/store"
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
	// initialize our redis client
	store.NewClient()

	// init the bot
	bot, err := discordgo.New("Bot " + viper.GetString("token"))
	if err != nil {
		log.Warn("error creating discord session: ", err)
	}

	bot.AddHandler(messageCreate)

	err = bot.Open()
	if err != nil {
		log.Error("error opening connection, ", err)
		return
	}

	bot.State.MaxMessageCount = 50

	router := exrouter.New()
	// Ping Pong
	//router.On("ping", ping.Route).Desc("sends a ping/pong")
	router.OnMatch("PingPong", dgrouter.NewRegexMatcher("p(i|o)ng"), ping.Route).Desc("sends a ping or pong")
	// Karma
	router.On("karma", karma.GetKarma).Desc("gets karma by user or your karma if no user specified\n\t\t" +
		"@user ++ will add karma\n\t\t" +
		"@user -- will remove karma")
	bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		match, err := regexp.MatchString(`(\-\-|\+\+)`, m.Message.Content)
		if err != nil {
			log.Error(err)
		} else if match {
			karma.Handler(s, m)
		}
	})
	// Quotes
	router.On("grab", quotes.Grab).Desc("grab quote by user or phrase")
	router.On("rq", quotes.RandomQuote).Desc("recall a random quote")

	// Help
	router.Default = router.On("help", func(ctx *exrouter.Context) {
		var text = ""
		for _, v := range router.Routes {
			text += v.Name + " : \t" + v.Description + "\n"
		}
		ctx.Reply("```" + text + "```")
	}).Desc("prints this help menu")

	// Establish routes for all commands managed by router
	bot.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(bot, viper.GetString("prefix"), bot.State.User.ID, m.Message)
	})

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
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
	messageLogger.Info(m.ContentWithMentionsReplaced())
}
