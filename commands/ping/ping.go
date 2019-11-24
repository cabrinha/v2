package ping

import (
	"github.com/Necroforger/dgrouter/exrouter"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Route sends a ping or pong back to Author
func Route(ctx *exrouter.Context) {
	cmd := ctx.Msg.ContentWithMentionsReplaced()
	if strings.Contains(cmd, "ping") {
		log.Info("Ping received from %s, sending pong.", ctx.Msg.Author.Username)
		ctx.Reply("Pong!")
	} else if strings.Contains(cmd, "pong") {
		log.Info("Pong received from %s, sending ping.", ctx.Msg.Author.Username)
		ctx.Reply("Ping!")
	}
}
