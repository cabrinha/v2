package karma

import (
	"fmt"
	"strings"
	"regexp"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/cabrinha/v2/plugins/store"
	"github.com/go-redis/redis"
)

func hasMentions(ctx *exrouter.Context) bool {
	if len(ctx.Msg.Mentions) > 0 {
		return true
	}
	return false
}

// GetKarma gets a user's karma score and returns it
func GetKarma(ctx *exrouter.Context) {
	if hasMentions(ctx) {
		var scores []string
		for _, u := range ctx.Msg.Mentions {
			scores = append(scores, u.Username, getScore(u))
		}
		if scores != nil {
			ctx.Reply(scores)
		} else {
			ctx.Reply("No karma found.")
		}
	} else {
		ctx.Reply(ctx.Msg.Author.Mention(), ": your score is: ", getScore(ctx.Msg.Author))
	}
}

// Need a new redis client here
var red = store.NewClient()

// Just get the score for a user
func getScore(user *discordgo.User) string {
	result, err := red.HGet(user.ID, "karma").Result()
	if err == redis.Nil {
		fmt.Println("error fetching karma score for user: ", user, err)
		fmt.Println("creating a 0 score for user: ", user)
		red.HSet(user.ID, "karma", 0)
		result, err = red.HGet(user.ID, "karma").Result()
	}
	return result
}

// Handler will handle karma commands
func Handler(ctx *exrouter.Context) {
	if strings.Contains(ctx.Msg.Content, fmt.Sprintf("%s .*\+\+.*\|.*\-\-.*", ctx.Msg.Author.ID) {
		ctx.Reply("You can't alter your own karma.")
	}
}
