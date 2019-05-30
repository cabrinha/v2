package quotes

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/cabrinha/v2/plugins/store"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// Grab a quote from known messages in a channel
func Grab(ctx *exrouter.Context) {
	args := ctx.Args.After(1)
	channel, _ := ctx.Ses.State.Channel(ctx.Msg.ChannelID)
	messages := channel.Messages
	log.Infof("Grab received from %s with args: %s", ctx.Msg.Author.Username, args)

	// If we have no args, grab the last message said
	if args == "" {
		// -2: we don't want to grab the `grab` command
		msg := messages[len(messages)-2]
		// If the bot said it, return
		if msg.Author.ID == ctx.Ses.State.User.ID {
			return
		}
		grabQuote(ctx, msg.Author, msg.ContentWithMentionsReplaced())
		return
	}

	// reverse the messages so we're always grabbing the newest message
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if strings.Contains(msg.Content, args) && !strings.Contains(msg.Content, "!grab") { // replace !grab with a regex later
			log.Infof("Found quote from %s: %s", msg.Author.Username, msg.ContentWithMentionsReplaced())
			grabQuote(ctx, msg.Author, msg.ContentWithMentionsReplaced())
			break
		}
	}
}

func grabQuote(ctx *exrouter.Context, user *discordgo.User, quote string) {
	result := store.Client.LPush(fmt.Sprintf("%s:quotes", user.ID), user.Username+": "+quote)
	if result.Err() != redis.Nil {
		log.Infof("Stored new quote for %s: %s", user.Username, quote)
		ctx.Reply("Grabbed.")
	} else {
		log.Errorf("Unable to store quote for %s: err: %v", user.Username, result.Err())
	}
}

// RandomQuote selects a quote at random and sends it to the channel
func RandomQuote(ctx *exrouter.Context) {
	log.Infof("RandomQuote received from %s with args: %s", ctx.Msg.Author.Username, ctx.Args.After(1))
	// init rand seed
	rand.Seed(time.Now().Unix())
	// get all quotes from redis
	allQuotes := store.Client.Keys("*:quotes").Val()
	if len(allQuotes) != 0 {
		k := rand.Int() % len(allQuotes)                          // select a key
		i := rand.Int63() % store.Client.LLen(allQuotes[k]).Val() // select an index
		randQuote := store.Client.LIndex(allQuotes[k], i).Val()
		ctx.Reply("```" + randQuote + "```")
	}
}
