package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var channelHourLimit = 1.0

func getChannelLogs(c echo.Context) error {
	channel := strings.TrimSpace(strings.ToLower(c.Param("channel")))

	fromTime, toTime, err := parseFromTo(c.QueryParam("from"), c.QueryParam("to"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	channelid := getUserid(channel)

	var logResult chatLog

	orderBy := orderAsc
	_, reverse := c.QueryParams()["reverse"]
	if reverse {
		orderBy = orderDesc
	}

	limit := c.QueryParam("limit")
	limitInt := 0
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)

		if err != nil || limitInt < 1 {
			return c.JSON(http.StatusBadRequest, "Invalid limit")
		}
	}

	selectFields := []string{"message", "timestamp", "userid"}
	whereClauses := []string{"channelid = ?", "timestamp >= ?", "timestamp <= ?"}

	iter := cassandra.Query(
		buildQuery(selectFields, "logstv.channel_messages", whereClauses, orderBy, limitInt),
		channelid,
		fromTime,
		toTime,
	).Iter()

	var message chatMessage
	var ts time.Time
	var fetchedUserid int64
	var messageRaw string
	for iter.Scan(&messageRaw, &ts, &fetchedUserid) {
		_, user, parsedMessage := twitch.ParseMessage(messageRaw)

		message.Timestamp = timestamp{ts}
		message.Username = user.Username
		message.Text = parsedMessage.Text
		message.Type = parsedMessage.Type

		logResult.Messages = append(logResult.Messages, message)
	}

	if err := iter.Close(); err != nil {
		log.Error(err)
	}

	if c.Request().Header.Get("Content-Type") == "application/json" || c.QueryParam("type") == "json" {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		return json.NewEncoder(c.Response()).Encode(logResult)
	}

	c.Response().WriteHeader(http.StatusOK)
	writeTextChatLog(&logResult, c.Response())

	return nil
}
