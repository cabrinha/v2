package karma

import (
	"fmt"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/cabrinha/v2/plugins/store"
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
			scores = append(scores, "%d: %s -", u.ID, getScore(u.ID))
		}
		if scores != nil {
			ctx.Reply(scores)
		} else {
			ctx.Reply("No karma found.")
		}
	}
	if getScore(ctx.Msg.Author.ID) != "" {
		ctx.Reply(getScore(ctx.Msg.Author.ID))
	} else {
		ctx.Reply("No karma found.")
	}
}

// Just get the score for a user
func getScore(user string) string {
	result, err := store.Client.HGet(user, "karma").Result()
	if err != nil {
		fmt.Println("error fetching karma score for user: ", user, err)
	}
	return result
}
