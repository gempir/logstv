package main

import (
	"log"

	"github.com/gempir/go-twitch-irc"
	"github.com/gempir/logstv/common"
	"github.com/gocql/gocql"
)

var cassandra *gocql.Session
var tClient *twitch.Client

func main() {
	common.LoadEnv()
	startup()
	defer cassandra.Close()

	tClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		go handleMessage(channel, user, message)
	})

	tClient.Join("gempir")

	// joinSavedChannels()

	// go func() {
	// 	for {
	// 		joinTop1000Channels()
	// 		joinSavedChannels()
	// 		time.Sleep(time.Minute * 15)
	// 	}
	// }()

	tClient.Connect()
}

func handleMessage(channel string, user twitch.User, message twitch.Message) {
	err := cassandra.Query("INSERT INTO streamlogs.messages (channelId, userId, message, timestamp) VALUES (?, ?, ?, ?)", message.Tags["room-id"], user.UserID, message.Text, message.Time).Exec()
	if err != nil {
		log.Printf("Failed message INSERT %s", err.Error())
	}

	err = cassandra.Query("INSERT INTO streamlogs.channels (userId, username) VALUES (?, ?) IF NOT EXISTS", user.UserID, user.Username).Exec()
	if err != nil {
		log.Printf("Failed channel INSERT %s", err.Error())
	}
}
