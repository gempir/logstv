package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gempir/streamlogs/common"
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

func joinTop1000Channels() {
	top := getTopChannels(0)
	top = append(top, getTopChannels(100)...)
	top = append(top, getTopChannels(200)...)
	top = append(top, getTopChannels(300)...)
	top = append(top, getTopChannels(400)...)
	top = append(top, getTopChannels(500)...)
	top = append(top, getTopChannels(600)...)
	top = append(top, getTopChannels(700)...)
	top = append(top, getTopChannels(800)...)
	top = append(top, getTopChannels(900)...)
	top = append(top, getTopChannels(1000)...)

	for _, channel := range top {
		joinChannel(channel.Channel.UserID, channel.Channel.Name)
	}
}

func joinSavedChannels() {
	var channelName string
	var channelID int64

	iter := cassandra.Query("SELECT userId,username FROM streamlogs.channels").Iter()
	for iter.Scan(&channelID, &channelName) {
		joinChannel(channelID, channelName)
	}
}

func joinChannel(channelID int64, channelName string) {
	err := cassandra.Query("INSERT INTO streamlogs.channels (userId, username) VALUES (?, ?) IF NOT EXISTS", channelID, channelName).Exec()
	if err != nil {
		fmt.Printf("Failed to insert channel: %s", err.Error())
	}

	if isJoinedChannel(channelName) {
		return
	}

	fmt.Printf("Joining: %s\r\n", channelName)
	tClient.Join(channelName)
	joinedChannels = append(joinedChannels, channelName)
}

func getTopChannels(offset int) []Stream {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/streams?limit=100&offset=%d", offset), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Client-Id", common.GetEnv("CLIENTID"))
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var streams Streams
	json.Unmarshal(contents, &streams)

	return streams.Streams
}

func isJoinedChannel(channelName string) bool {
	for _, username := range joinedChannels {
		if username == channelName {
			return true
		}
	}
	return false
}
