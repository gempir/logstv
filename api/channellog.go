package main

import (
	"net/http"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var channelHourLimit = 1.0

func getChannelLogs(c echo.Context) error {
	channel := strings.TrimSpace(strings.ToLower(c.Param("channel")))

	fromTime, toTime, err := parseFromTo(c.QueryParam("from"), c.QueryParam("to"), channelHourLimit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	channelid := getUserid(channel)

	var logResult chatLog

	orderBy := buildOrder(c)
	limit, err := buildLimit(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid limit")
	}

	selectFields := []string{"message", "dateOf(timeuuid)"}
	whereClauses := []string{"channelid = ?", "timeuuid >= maxTimeuuid(?)", "timeuuid <= minTimeuuid(?)"}

	iter := cassandra.Query(
		buildQuery(selectFields, "logstv.channel_messages", whereClauses, orderBy, limit),
		channelid,
		fromTime,
		toTime,
	).Iter()

	var message chatMessage
	var ts time.Time
	var messageRaw string
	for iter.Scan(&messageRaw, &ts) {
		channel, user, parsedMessage := twitch.ParseMessage(messageRaw)

		message.Timestamp = timestamp{ts}
		message.Username = user.Username
		message.Text = parsedMessage.Text
		message.Type = parsedMessage.Type
		message.Channel = channel

		logResult.Messages = append(logResult.Messages, message)
	}

	if err := iter.Close(); err != nil {
		log.Error(err)
	}

	if c.Request().Header.Get("Content-Type") == "application/json" || c.QueryParam("type") == "json" {
		return writeJSONResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}
