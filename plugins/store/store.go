package store

import (
	"github.com/go-redis/redis"
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
}
