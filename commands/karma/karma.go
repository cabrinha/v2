package karma

import (
	"fmt"
	"strconv"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/cabrinha/v2/plugins/store"
)

// GetKarma gets a user's karma score and returns it
func GetKarma(ctx *exrouter.Context) {
	ctx.Reply(getScore(ctx))
}

// Just get the score for a user
func getScore(ctx *exrouter.Context) int {
	user := ctx.Msg.Author.String()
	result, err := store.Client.HGet(user, "karma").Result()
	if err != nil {
		fmt.Println("error fetching karma score for user: ", user, err)
	}
	score, _ := strconv.Atoi(result)

	return score
}

// Plus adds 1 to user's karma score
func Plus(user string, ctx *exrouter.Context) {
	if user == ctx.Msg.Author.String() {
		ctx.Reply("You cannot alter your own karma.")
		return
	}
	newScore := getScore(ctx) + 1
	if store.Client.HSet(user, "karma", newScore).Val() == true {
		ctx.Reply(fmt.Printf("%s karma increased to: %d", user, newScore))
	} else {
		fmt.Println("Failed to set karma on user: ", user)
	}
}

// Minus subtracts 1 to user's karma score
func Minus(user string, ctx *exrouter.Context) {
	if user == ctx.Msg.Author.String() {
		ctx.Reply("You cannot alter your own karma.")
		return
	}
	newScore := getScore(ctx) - 1
	if store.Client.HSet(user, "karma", newScore).Val() == true {
		ctx.Reply(fmt.Printf("%s karma decreased to: %d", user, newScore))
	} else {
		fmt.Println("Failed to set karma on user: ", user)
	}
}
