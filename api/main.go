package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	minio "github.com/minio/minio-go"

	twitch "github.com/gempir/go-twitch-irc"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	s3Client *minio.Client
	s3Bucket string
)

func main() {
	var err error
	s3Client, err = minio.NewV2(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_ID"), os.Getenv("S3_ACCESS_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	s3Bucket = os.Getenv("S3_BUCKET")

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
	e.GET("/channelid/:channelid/userid/:userid/:year/:month", getChannelUserLogs)

	fmt.Println("starting streamlogs API on port :8010")
	log.Fatal(e.Start(":8010"))
}

type chatLog struct {
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Text      string             `json:"text"`
	Username  string             `json:"username"`
	Channel   string             `json:"channel"`
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

func writeTextResponse(c echo.Context, cLog *chatLog) error {
	c.Response().WriteHeader(http.StatusOK)

	for _, cMessage := range cLog.Messages {
		switch cMessage.Type {
		case twitch.PRIVMSG:
			c.Response().Write([]byte(fmt.Sprintf("[%s] #%s %s: %s\r\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Username, cMessage.Text)))
			break
		case twitch.CLEARCHAT:
			c.Response().Write([]byte(fmt.Sprintf("[%s] #%s %s\r\n", cMessage.Timestamp.Format("2006-01-2 15:04:05"), cMessage.Channel, cMessage.Text)))
			break
		}
	}

	return nil
}

func writeJSONResponse(c echo.Context, logResult *chatLog) error {
	_, stream := c.QueryParams()["stream"]
	if stream {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		return json.NewEncoder(c.Response()).Encode(logResult)
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(logResult)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, echo.MIMEApplicationJSONCharsetUTF8, data)
}
