package ping

import "github.com/Necroforger/dgrouter/exrouter"

func PingRoute(ctx *exrouter.Context) {
	ctx.Reply("Pong!")
}

func PongRoute(ctx *exrouter.Context) {
	ctx.Reply("Ping!")
}
