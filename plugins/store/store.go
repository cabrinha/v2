package store

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Client for reuse elsewhere
var Client *redis.Client

// NewClient creates a redis client connection
func NewClient() {
	Client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis_host"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := Client.Ping().Result()
	log.WithFields(log.Fields{
		"PING":  pong,
		"Error": err,
	}).Info("Connecting to redis")
}
