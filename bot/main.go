package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
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
	Name string `json:"name"`
}

var keyspaceQuery = `
	CREATE  KEYSPACE IF NOT EXISTS streamlogs
	WITH REPLICATION = { 
		'class' : 'SimpleStrategy', 
		'replication_factor' : 1 
	};`

var tableQuery = `
CREATE TABLE IF NOT EXISTS streamlogs.messages (
	id uuid,
	channelId bigint,
	userId bigint,
	message text,
	timestamp timestamp,
	PRIMARY KEY (id)
);`

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	tclient := twitch.NewClient("justinfan123123", "oauth:123123123")

	hosts := strings.Split(os.Getenv("DBHOSTS"), ",")
	cluster := gocql.NewCluster(hosts...)

	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	err = session.Query(keyspaceQuery).Exec()
	if err != nil {
		panic(err)
	}

	err = session.Query(tableQuery).Exec()
	if err != nil {
		panic(err)
	}

	tclient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		fmt.Println(message.Text)
		err = session.Query("INSERT INTO streamlogs.messages (id, channelId, userId, message, timestamp) VALUES (?, ?, ?, ?, ?)", message.Tags["id"], message.Tags["room-id"], user.UserID, message.Text, message.Time).Exec()
		if err != nil {
			log.Printf("Failed INSERT %s", err.Error())
		}
	})

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
	top = append(top, getTopChannels(1100)...)
	top = append(top, getTopChannels(1200)...)

	fmt.Println(top)
	for _, channel := range top {
		fmt.Printf("Joining: %s\r\n", channel.Channel.Name)
		go tclient.Join(channel.Channel.Name)
	}

	// go func() {
	// 	for {
	// 		top := getTopChannels()
	// 		for _, channel := range top {
	// 			fmt.Printf("Joining: %s\r\n", channel.Channel.Name)
	// 			go tclient.Join(channel.Channel.Name)
	// 		}
	// 		time.Sleep(time.Hour)
	// 	}
	// }()

	go tclient.Join("pajlada")

	tclient.Connect()
}

func getTopChannels(offset int) []Stream {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/streams?limit=100&offset=%d", offset), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Client-Id", getEnv("CLIENTID"))
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(contents))
	var streams Streams
	json.Unmarshal(contents, &streams)

	return streams.Streams
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Missing env var: " + key)
}
