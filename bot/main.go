package main

import (
	"time"

	"github.com/gempir/go-twitch-irc"
	"github.com/gempir/logstv/common"
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

var cassandra *gocql.Session
var tClient *twitch.Client
var userCache = make(map[int64]bool)

func main() {
	common.LoadEnv()
	startup()

	tClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		go handleMessage(channel, user, message)
	})

	go func() {
		for {
			joinSavedChannels()
			time.Sleep(time.Minute)
		}
	}()

	panic(tClient.Connect())
}

func handleMessage(channel string, user twitch.User, message twitch.Message) {
	go func() {
		err := cassandra.Query("INSERT INTO logstv.messages (channelId, userId, message, timestamp) VALUES (?, ?, ?, ?)", message.Tags["room-id"], user.UserID, message.Text, message.Time).Exec()
		if err != nil {
			log.Errorf("Failed message INSERT %s", err.Error())
		}
	}()
	go func() {
		if _, ok := userCache[user.UserID]; ok {
			return
		}

		userCache[user.UserID] = true
		err := cassandra.Query("INSERT INTO logstv.users (userId, username) VALUES (?, ?) IF NOT EXISTS", user.UserID, user.Username).Exec()
		if err != nil {
			log.Errorf("Failed channel INSERT %s", err.Error())
		}
	}()
}
