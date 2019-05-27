package karma

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/cabrinha/v2/plugins/store"
	"github.com/go-redis/redis"
)

// Need a new redis client here
var redisdb = store.NewClient()

// Translate a User ID to a Username
func userNameFromID(ctx *exrouter.Context, u string) string {
	member, err := ctx.Ses.State.Member(ctx.Msg.GuildID, u)
	if err != nil {
		log.Errorf("Fetching member %s failed: %s", u, err)
	}
	name := strings.Split(member.User.Username, "#")[0]
	return fmt.Sprint(name)
}

// Check if the message has mentions
func hasMentions(ctx *exrouter.Context) bool {
	if len(ctx.Msg.Mentions) > 0 {
		return true
	}
	return false
}

// GetKarma gets a user's karma score and returns it
func GetKarma(ctx *exrouter.Context) {
	if hasMentions(ctx) {
		scores := make(map[string]string)
		for _, u := range ctx.Msg.Mentions {
			score, err := getScore(ctx, u.ID)
			if err != nil {
				log.Error(err)
			}
			scores[u.Username] = strconv.Itoa(score)
		}
		for k, v := range scores {
			ctx.Reply(fmt.Sprintf("%s: %s", k, v))
		}
	} else {
		score, err := getScore(ctx, ctx.Msg.Author.ID)
		if err != nil {
			log.Error(err)
		}
		ctx.Reply(fmt.Sprintf("%s, your score is %d", ctx.Msg.Author.Mention(), score))
	}
}

// If the author's User ID is in the content, slap their hand
func authorInContent(ctx *exrouter.Context) bool {
	if strings.Contains(ctx.Msg.Content, ctx.Msg.Author.ID) {
		ctx.Reply("You can't alter your own karma.")
		return true
	}
	return false
}

// Just get the score for a user
func getScore(ctx *exrouter.Context, user string) (int, error) {
	log.Infof("Getting karma score for user: %s", userNameFromID(ctx, user))
	result, err := redisdb.HGet(user, "karma").Result()
	if err == redis.Nil {
		log.Error("error fetching karma score for user: ", user, err)
		log.Info("creating a 0 score for user: ", user)
		redisdb.HSet(user, "karma", 0)
		return 0, nil
	}
	return strconv.Atoi(result)
}

func plus(ctx *exrouter.Context, u string) int {
	var i int
	for _, u := range ctx.Msg.Mentions {
		i, _ = getScore(ctx, u.ID)
		newScore := i + 1
		result := redisdb.HSet(u.ID, "karma", newScore)
		if result.Err() != redis.Nil {
			log.Infof("Set new score for user ID: %s to %d", u, newScore)
			return newScore
		}
		log.Errorf("Unable to set score for user ID: %s, Err: %s", u, result.Err())
		ctx.Reply(fmt.Sprintf("Unable to alter the karma score for: %s", userNameFromID(ctx, u.ID)))
	}
	return i
}

func minus(ctx *exrouter.Context, u string) int {
	var i int
	for _, u := range ctx.Msg.Mentions {
		i, _ = getScore(ctx, u.ID)
		newScore := i - 1
		result := redisdb.HSet(u.ID, "karma", newScore)
		if result.Err() != redis.Nil {
			log.Infof("Set new score for user ID: %s to %d", u, newScore)
			return newScore
		}
		log.Errorf("Unable to set score for user ID: %s, Err: %s", u, result.Err())
		ctx.Reply(fmt.Sprintf("Unable to alter the karma score for: %s", userNameFromID(ctx, u.ID)))
	}
	return i
}

// Handler will handle karma commands
func Handler(ctx *exrouter.Context) {
	rePlus := regexp.MustCompile(`\!\+\+\s+<@(?P<userID>\d{18}).*`)
	reMinus := regexp.MustCompile(`\!\-\-\s+<@(?P<userID>\d{18}).*`)
	if !authorInContent(ctx) {
		plusMatch := rePlus.FindStringSubmatch(ctx.Msg.Content)
		if len(plusMatch) > 0 {
			result := make(map[string]string)
			for i, name := range rePlus.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = plusMatch[i]
				}
			}
			ctx.Reply(fmt.Sprintf("%s's karma is now at: %d", userNameFromID(ctx, result["userID"]), plus(ctx, result["userID"])))
		} else {
			minusMatch := reMinus.FindStringSubmatch(ctx.Msg.Content)
			if len(minusMatch) > 0 {
				result := make(map[string]string)
				for i, name := range reMinus.SubexpNames() {
					if i != 0 && name != "" {
						result[name] = minusMatch[i]
					}
				}
				ctx.Reply(
					fmt.Sprintf("%s's karma is now at: %d", userNameFromID(ctx, result["userID"]), minus(ctx, result["userID"])))
			}
		}
	}
}
