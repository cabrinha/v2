package karma

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/cabrinha/v2/plugins/store"
	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
)

// Translate a User ID to a Username
func userNameFromID(s *discordgo.Session, m *discordgo.Message, u string) string {
	member, err := s.State.Member(m.GuildID, u)
	if err != nil {
		log.Errorf("Fetching member %s failed: %s", u, err)
	}
	name := strings.Split(member.User.Username, "#")[0]
	return fmt.Sprint(name)
}

// Check if the message has mentions
func hasMentions(m *discordgo.Message) bool {
	return len(m.Mentions) > 0
}

// Just get the score for a user
func getScore(s *discordgo.Session, m *discordgo.Message, u string) (int, error) {
	log.Infof("Getting karma score for user: %s", userNameFromID(s, m, u))
	result, err := store.Client.Get(fmt.Sprintf("%s:karma", u)).Result()
	if err == redis.Nil {
		log.Error("Error fetching karma score for user: ", u, err)
		log.Info("Creating a 0 score for user: ", u)
		store.Client.Set(fmt.Sprintf("%s:karma", u), 0, 0)
		return 0, nil
	}
	return strconv.Atoi(result)
}

// Alter the karma for a user
func plus(s *discordgo.Session, m *discordgo.Message) int {
	var i int
	for _, u := range m.Mentions {
		i, _ = getScore(s, m, u.ID)
		newScore := i + 1
		result := store.Client.Set(fmt.Sprintf("%s:karma", u.ID), newScore, 0)
		if result.Err() != redis.Nil {
			log.Infof("Set new score for user ID: %s to %d", u, newScore)
			return newScore
		}
		log.Errorf("Unable to set score for user ID: %s, Err: %s", u, result.Err())
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to alter the karma score for: %s", userNameFromID(s, m, u.ID)))
	}
	return i
}

func minus(s *discordgo.Session, m *discordgo.Message) int {
	var i int
	for _, u := range m.Mentions {
		i, _ = getScore(s, m, u.ID)
		newScore := i - 1
		result := store.Client.Set(fmt.Sprintf("%s:karma", u), newScore, 0)
		if result.Err() != redis.Nil {
			log.Infof("Set new score for user ID: %s to %d", u, newScore)
			return newScore
		}
		log.Errorf("Unable to set score for user ID: %s, Err: %s", u, result.Err())
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to alter the karma score for: %s", userNameFromID(s, m, u.ID)))
	}
	return i
}

// Check if we have a plus or minus match, return bool
func isPlus(message string) bool {
	plusRegex := regexp.MustCompile(`(.*)?<@(!)?(?P<userID>\d{18})>\s+\+\+(.*)?`)
	matched := plusRegex.MatchString(message)
	return matched
}

func isMinus(message string) bool {
	minusRegex := regexp.MustCompile(`(.*)?<@(!)?(?P<userID>\d{18})>\s+--(.*)?`)
	matched := minusRegex.MatchString(message)
	return matched
}

// Alter karma in a given direction based on regex
func alterKarma(p string, m *discordgo.MessageCreate, s *discordgo.Session, plusOrMinus func(s *discordgo.Session, m *discordgo.Message) int) {
	re := regexp.MustCompile(p)
	match := re.FindStringSubmatch(m.Message.Content)
	if len(match) > 0 {
		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		s.ChannelMessageSend(m.Message.ChannelID, fmt.Sprintf(
			"%s's karma is now at: %d",
			userNameFromID(s, m.Message, result["userID"]),
			plusOrMinus(s, m.Message),
		))
	}
}

// GetKarma gets a user's karma score and returns it
func GetKarma(ctx *exrouter.Context) {
	if hasMentions(ctx.Msg) {
		scores := make(map[string]string)
		for _, u := range ctx.Msg.Mentions {
			score, err := getScore(ctx.Ses, ctx.Msg, u.ID)
			if err != nil {
				log.Error(err)
			}
			scores[u.Username] = strconv.Itoa(score)
		}
		for k, v := range scores {
			ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf("%s: %s", k, v))
		}
	} else {
		score, err := getScore(ctx.Ses, ctx.Msg, ctx.Msg.Author.ID)
		if err != nil {
			log.Error(err)
		}
		ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf("%s, your score is %d", ctx.Msg.Author.Mention(), score))
	}
}

// Handler handles the updating of karma scores
func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Message.Author.ID == s.State.User.ID {
		return
	}

	// If the author is trying to message their own karma: don't
	for _, u := range m.Message.Mentions {
		if u.ID == m.Message.Author.ID {
			s.ChannelMessageSend(m.Message.ChannelID, "You can't alter your own karma.")
			return
		}
	}

	if isPlus(m.Message.Content) {
		alterKarma(`(.*)?<@(!)?(?P<userID>\d{18})>\s+\+\+(.*)?`, m, s, plus)
	}
	if isMinus(m.Message.Content) {
		alterKarma(`(.*)?<@(!)?(?P<userID>\d{18})>\s+--(.*)?`, m, s, minus)
	}
}
