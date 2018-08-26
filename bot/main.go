package main

import (
	"github.com/gempir/go-twitch-irc"
	"github.com/gempir/logstv/common"
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

var cassandra *gocql.Session
var tClient *twitch.Client

func main() {
	common.LoadEnv()
	startup()

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
	go func() {
		err := cassandra.Query("INSERT INTO logstv.messages (channelId, userId, message, timestamp) VALUES (?, ?, ?, ?)", message.Tags["room-id"], user.UserID, message.Text, message.Time).Exec()
		if err != nil {
			log.Errorf("Failed message INSERT %s", err.Error())
		}
	}()
	go func() {
		err := cassandra.Query("INSERT INTO logstv.channels (userId, username) VALUES (?, ?) IF NOT EXISTS", user.UserID, user.Username).Exec()
		if err != nil {
			log.Errorf("Failed channel INSERT %s", err.Error())
		}
	}()
}
