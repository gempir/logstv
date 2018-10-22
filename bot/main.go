package main

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"
	"github.com/go-redis/redis"
)

var (
	rClient *redis.Client
)

func main() {
	fmt.Println("Bot booting up")
	tClient := twitch.NewClient("justinfan123123", "oauth:123123123")

	rClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rClient.Ping().Result()
	if err != nil {
		panic("Redis unavailable: " + err.Error())
	}

	tClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		go rClient.LPush("messages", message.Raw)
	})

	tClient.OnNewClearchatMessage(func(channel string, user twitch.User, message twitch.Message) {
		go rClient.LPush("messages", message.Raw)
	})

	tClient.OnNewUsernoticeMessage(func(channel string, user twitch.User, message twitch.Message) {
		go rClient.LPush("messages", message.Raw)
	})

	tClient.Join("gempir")
	// tClient.Join("pajlada")
	// tClient.Join("nymn")
	// tClient.Join("forsen")

	panic(tClient.Connect())
}
