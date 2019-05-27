package karma

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"

	"github.com/cabrinha/v2/plugins/store"
	"github.com/go-redis/redis"
)

// Need a new redis client here
var redisdb = store.NewClient()

// Translate a User ID to a Username
func userNameFromID(s *discordgo.Session, m *discordgo.MessageCreate, u string) string {
	member, err := s.State.Member(m.GuildID, u)
	if err != nil {
		log.Errorf("Fetching member %s failed: %s", u, err)
	}
	name := strings.Split(member.User.Username, "#")[0]
	return fmt.Sprint(name)
}

// Check if the message has mentions
func hasMentions(m *discordgo.MessageCreate) bool {
	if len(m.Mentions) > 0 {
		return true
	}
	return false
}

// Just get the score for a user
func getScore(s *discordgo.Session, m *discordgo.MessageCreate, u string) (int, error) {
	log.Infof("Getting karma score for user: %s", userNameFromID(s, m, u))
	result, err := redisdb.HGet(u, "karma").Result()
	if err == redis.Nil {
		log.Error("error fetching karma score for user: ", u, err)
		log.Info("creating a 0 score for user: ", u)
		redisdb.HSet(u, "karma", 0)
		return 0, nil
	}
	return strconv.Atoi(result)
}

func plus(s *discordgo.Session, m *discordgo.MessageCreate, u string) int {
	var i int
	for _, u := range m.Mentions {
		i, _ = getScore(s, m, u.ID)
		newScore := i + 1
		result := redisdb.HSet(u.ID, "karma", newScore)
		if result.Err() != redis.Nil {
			log.Infof("Set new score for user ID: %s to %d", u, newScore)
			return newScore
		}
		log.Errorf("Unable to set score for user ID: %s, Err: %s", u, result.Err())
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to alter the karma score for: %s", userNameFromID(s, m, u.ID)))
	}
	return i
}

func minus(s *discordgo.Session, m *discordgo.MessageCreate, u string) int {
	var i int
	for _, u := range m.Mentions {
		i, _ = getScore(s, m, u.ID)
		newScore := i - 1
		result := redisdb.HSet(u.ID, "karma", newScore)
		if result.Err() != redis.Nil {
			log.Infof("Set new score for user ID: %s to %d", u, newScore)
			return newScore
		}
		log.Errorf("Unable to set score for user ID: %s, Err: %s", u, result.Err())
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to alter the karma score for: %s", userNameFromID(s, m, u.ID)))
	}
	return i
}

// Check if we have a match, return bool and capture group names
func isPlus(user string, message string) bool {
	plusRegex := regexp.MustCompile(`(.*)?<@(!)?(?P<userID>\d{18})>\s+\+\+(.*)?`)
	matched, err := regexp.MatchString(fmt.Sprintf(plusRegex.String(), user), message)
	if err != nil {
		panic(err)
	}
	return matched
}

func isMinus(user string, message string) bool {
	minusRegex := regexp.MustCompile(`(.*)?<@(!)?(?P<userID>\d{18})>\s+--(.*)?`)
	matched, err := regexp.MatchString(fmt.Sprintf(minusRegex.String(), user), message)
	if err != nil {
		panic(err)
	}
	return matched
}

// GetKarma gets a user's karma score and returns it
func GetKarma(s *discordgo.Session, m *discordgo.MessageCreate) {
	if hasMentions(m) {
		scores := make(map[string]string)
		for _, u := range m.Mentions {
			score, err := getScore(s, m, u.ID)
			if err != nil {
				log.Error(err)
			}
			scores[u.Username] = strconv.Itoa(score)
		}
		for k, v := range scores {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", k, v))
		}
	} else {
		score, err := getScore(s, m, m.Author.ID)
		if err != nil {
			log.Error(err)
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, your score is %d", m.Author.Mention(), score))
	}
}

// Handler handles the updating of karma scores
func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the author is trying to message their own karma: don't
	for _, u := range m.Mentions {
		if u.ID == m.Author.ID {
			s.ChannelMessageSend(m.ChannelID, "You can't alter your own karma.")
			return
		}
	}

	for _, u := range m.Mentions {
		if isPlus(u.ID, m.Content) {
			plusRegex := regexp.MustCompile(`(.*)?<@(!)?(?P<userID>\d{18})>\s+\+\+(.*)?`)
			plusMatch := plusRegex.FindStringSubmatch(m.Content)
			if len(plusMatch) > 0 {
				result := make(map[string]string)
				for i, name := range plusRegex.SubexpNames() {
					if i != 0 && name != "" {
						result[name] = plusMatch[i]
					}
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's karma is now at: %d", userNameFromID(s, m, result["userID"]), plus(s, m, result["userID"])))
				}
			}
		}
		if isMinus(u.ID, m.Content) {
			minusRegex := regexp.MustCompile(`(.*)?<@(!)?(?P<userID>\d{18})>\s+--(.*)?`)
			minusMatch := minusRegex.FindStringSubmatch(m.Content)
			if len(minusMatch) > 0 {
				result := make(map[string]string)
				for i, name := range minusRegex.SubexpNames() {
					if i != 0 && name != "" {
						result[name] = minusMatch[i]
					}
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's karma is now at: %d", userNameFromID(s, m, result["userID"]), minus(s, m, result["userID"])))
				}
			}
		}
	}
}

/*
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
*/
