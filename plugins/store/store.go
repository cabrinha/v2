package store

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// NewClient creates a redis client connection
func NewClient() *redis.Client {
	Client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis_host"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return Client
}
