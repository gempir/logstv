package main

import (
	"strconv"
	"time"

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
		go persistUser(user.Username, user.UserID)

		channelid, err := strconv.ParseInt(message.Tags["room-id"], 10, 64)
		if err != nil {
			log.Errorf("Error parsing room-id to int64: %s", err.Error())
		}
		go persistMessage(channelid, user.UserID, message.Time, message.Raw)
	})

	tClient.OnNewClearchatMessage(func(channel string, user twitch.User, message twitch.Message) {
		userid, err := strconv.ParseInt(message.Tags["target-user-id"], 10, 64)
		if err != nil {
			log.Errorf("Error parsing target-user-id to int64: %s", err.Error())
		}

		go persistUser(user.Username, userid)

		channelid, err := strconv.ParseInt(message.Tags["room-id"], 10, 64)
		if err != nil {
			log.Errorf("Error parsing room-id to int64: %s", err.Error())
		}
		go persistMessage(channelid, userid, message.Time, message.Raw)
	})

	go func() {
		for {
			joinSavedChannels()
			time.Sleep(time.Minute)
		}
	}()

	panic(tClient.Connect())
}

func persistMessage(channelid int64, userid int64, time time.Time, messageRaw string) {
	go func() {
		err := cassandra.Query("INSERT INTO logstv.messages (channelid, userid, message, timeuuid) VALUES (?, ?, ?, ?)", channelid, userid, messageRaw, gocql.UUIDFromTime(time)).Exec()
		if err != nil {
			log.Errorf("Failed message INSERT %s", err.Error())
		}
	}()
	go func() {
		err := cassandra.Query("INSERT INTO logstv.channel_messages (channelid, message, timeuuid) VALUES (?, ?, ?)", channelid, messageRaw, gocql.UUIDFromTime(time)).Exec()
		if err != nil {
			log.Errorf("Failed channel_message INSERT %s", err.Error())
		}
	}()
	go func() {
		err := cassandra.Query("INSERT INTO logstv.user_messages (userid, message, timeuuid) VALUES (?, ?, ?)", userid, messageRaw, gocql.UUIDFromTime(time)).Exec()
		if err != nil {
			log.Errorf("Failed channel_message INSERT %s", err.Error())
		}
	}()
}

func persistUser(username string, userid int64) {
	go func() {
		err := cassandra.Query("INSERT INTO logstv.users (userid, username) VALUES (?, ?) IF NOT EXISTS", userid, username).Exec()
		if err != nil {
			log.Errorf("Failed channel INSERT %s", err.Error())
		}
	}()
}
