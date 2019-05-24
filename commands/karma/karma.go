package karma

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/cabrinha/v2/plugins/store"
)

// GetKarma gets a user's karma score and returns it
func GetKarma(ctx *exrouter.Context, user *discordgo.User) {
	ctx.Reply(getScore(ctx, user))
}

// Just get the score for a user
func getScore(ctx *exrouter.Context, user *discordgo.User) int {
	result, err := store.Client.HGet(user.ID, "karma").Result()
	if err != nil {
		fmt.Println("error fetching karma score for user: ", user, err)
	}
	score, _ := strconv.Atoi(result)

	return score
}

// IncKarma increments a user's karma by 1
func IncKarma(ctx *exrouter.Context, user *discordgo.User) {
	if user.ID == ctx.Msg.Author.ID {
		ctx.Reply("You cannot alter your own karma.")
	}
	newScore := getScore(ctx, user) + 1
	if store.Client.HSet(user.ID, "karma", newScore).Val() {
		ctx.Reply("Score for %s updated to: %d", user.String(), newScore)
	} else {
		fmt.Println("Failed to set karma on user: ", user.String())
	}
}

// DecKarma decrements a user's karma by 1
func DecKarma(ctx *exrouter.Context, user *discordgo.User) {
	if user.ID == ctx.Msg.Author.ID {
		ctx.Reply("You cannot alter your own karma.")
	}
	newScore := getScore(ctx, user) - 1
	if store.Client.HSet(user.ID, "karma", newScore).Val() {
		ctx.Reply("Score for %s updated to: %d", user.String(), newScore)
	} else {
		fmt.Println("Failed to set karma on user: ", user.String())
	}
}

// MentionsWithSuffix does ... idk
func MentionsWithSuffix(ctx *exrouter.Context, suffix string) []*discordgo.User {
	var matches []*discordgo.User
	i := strings.Index(ctx.Msg.Content, suffix)
	for i >= 0 {
		for _, user := range ctx.Msg.Mentions {
			if strings.HasSuffix(ctx.Msg.Content[:i-1], user.ID) {
				matches = append(matches, user)
			}
		}
		content := ctx.Msg.Content[i+len(suffix):]
		i = strings.Index(content, suffix)
	}
	return matches
}

// ApplyWithSuffix does something...
func ApplyWithSuffix(act func(ctx *exrouter.Context, user *discordgo.User), suffix string) func(*exrouter.Context) {
	return func(ctx *exrouter.Context) {
		matches := MentionsWithSuffix(ctx, suffix)
		for _, match := range matches {
			act(ctx, user)
		}
	}
}

// Plus adds 1 to mentioned user's karma score
func Plus(ctx *exrouter.Context, user *discordgo.User) {
	if user.ID != ctx.Msg.Author.ID {
		IncKarma(ctx, user)
	}
}

// Minus subtracts 1 from mentioned user's karma score
func Minus(ctx *exrouter.Context, user *discordgo.User) {
	if user.ID != ctx.Msg.Author.ID {
		DecKarma(ctx, user)
	}
}
