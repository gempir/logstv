package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocql/gocql"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/gempir/logstv/common"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var cassandra *gocql.Session

func main() {
	common.LoadEnv()

	var err error
	cassandra, err = common.NewDatabaseSession(strings.Split(common.GetEnv("DBHOSTS"), ","))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Debug = true

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to api.logs.tv!")
	})
	e.GET("/channel/:channel/user/:username", getUserLogs)
	e.GET("/channel/:channel", getChannelLogs)

	fmt.Println("starting streamlogs API on port :8010")
	log.Fatal(e.Start(":8010"))
}

type chatLog struct {
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Text      string             `json:"text"`
	Username  string             `json:"username"`
	Timestamp timestamp          `json:"timestamp"`
	Type      twitch.MessageType `json:"type"`
}

type timestamp struct {
	time.Time
}

func (t timestamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.UTC().Format(time.RFC3339) + "\""), nil
}

func (t *timestamp) UnmarshalJSON(data []byte) error {
	goTime, err := time.Parse(time.RFC3339, strings.TrimSuffix(strings.TrimPrefix(string(data[:]), "\""), "\""))
	if err != nil {
		return err
	}
	*t = timestamp{
		goTime,
	}
	return nil
}

func buildTextChatLog(cLog chatLog) string {
	var text string

	for _, cMessage := range cLog.Messages {
		switch cMessage.Type {
		case twitch.PRIVMSG:
			text += fmt.Sprintf("[%s] %s: %s\r\n", cMessage.Timestamp.Format("2006-01-2 15:04:05 UTC"), cMessage.Username, cMessage.Text)
			break
		case twitch.CLEARCHAT:
			text += fmt.Sprintf("[%s] %s\r\n", cMessage.Timestamp.Format("2006-01-2 15:04:05 UTC"), cMessage.Text)
			break
		}
	}

	return text
}
