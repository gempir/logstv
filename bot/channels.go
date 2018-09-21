package main

import (
	"fmt"
	"time"
)

type msg struct {
	Text     string    `json:"text"`
	Username string    `json:"username"`
	Time     time.Time `json:"time"`
}

// Streams struct
type Streams struct {
	Streams []Stream `json:"streams"`
}

// Stream struct
type Stream struct {
	Channel Channel `json:"channel"`
}

// Channel struct
type Channel struct {
	Name   string `json:"name"`
	UserID int64  `json:"_id"`
}

var joinedChannels = []string{}

func joinSavedChannels() {
	var channelName string

	iter := cassandra.Query("SELECT username FROM logstv.channels").Iter()
	for iter.Scan(&channelName) {
		joinChannel(channelName)
	}
}

func joinChannel(channelName string) {
	if isJoinedChannel(channelName) {
		return
	}

	fmt.Printf("Joining: %s\r\n", channelName)
	tClient.Join(channelName)
	joinedChannels = append(joinedChannels, channelName)
}

func isJoinedChannel(channelName string) bool {
	for _, username := range joinedChannels {
		if username == channelName {
			return true
		}
	}
	return false
}
