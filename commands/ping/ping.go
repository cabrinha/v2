package ping

import (
	"github.com/Necroforger/dgrouter/exrouter"
	log "github.com/sirupsen/logrus"
)

// PingRoute sends a ping back
func PingRoute(ctx *exrouter.Context) {
	log.Info("Ping received from %s, sending pong.", ctx.Msg.Author.Username)
	ctx.Reply("Pong!")
}

// PongRoute sends a ping back
func PongRoute(ctx *exrouter.Context) {
	log.Info("Ping received from %s, sending pong.", ctx.Msg.Author.Username)
	ctx.Reply("Ping!")
}
