package main

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/gempir/go-twitch-irc"
	minio "github.com/minio/minio-go"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var channelUserHourLimit = 744.0

func getChannelUserLogs(c echo.Context) error {
	channelid := c.Param("channelid")
	userid := c.Param("userid")
	year := c.Param("year")
	month := c.Param("month")

	objectName := fmt.Sprintf("%s_%s_%s", userid, year, month)

	object, err := s3Client.GetObject(s3Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Failure finding logs")
	}
	defer object.Close()

	scanner := bufio.NewScanner(object)

	var logResult chatLog
	for scanner.Scan() {
		line := scanner.Text()

		// first line is file name
		if line == objectName {
			continue
		}
		parsedChannel, user, parsedMessage := twitch.ParseMessage(line)
		if parsedMessage.Tags["room-id"] != channelid {
			continue
		}

		var msg chatMessage
		msg.Timestamp = timestamp{parsedMessage.Time}
		msg.Username = user.Username
		msg.Text = parsedMessage.Text
		msg.Type = parsedMessage.Type
		msg.Channel = parsedChannel

		logResult.Messages = append(logResult.Messages, msg)
	}

	if err := scanner.Err(); err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, "Failure reading messages")
	}

	if c.Request().Header.Get("Content-Type") == "application/json" || c.QueryParam("type") == "json" {
		return writeJSONResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}
