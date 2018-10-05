package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
)

func main() {
	twitchClient := twitch.NewClient("justinfan123123123", "oauth:123123123")
	fileLogger := NewFileLogger()

	go func() {
		for {
			data, err := ioutil.ReadFile("/etc/channels")
			if err != nil {
				panic(err)
			}

			channels := strings.Split(string(data), "\n")
			for _, channel := range channels {
				twitchClient.Join(channel)
			}

			time.Sleep(time.Minute)
		}
	}()

	twitchClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {

		go func() {
			err := fileLogger.LogMessageForUser(channel, user, message)
			if err != nil {
				log.Println(err.Error())
			}
		}()

		go func() {
			err := fileLogger.LogMessageForChannel(channel, user, message)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	})

	twitchClient.Connect()
}
