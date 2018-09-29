package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var userHourLimit = 744.0

func getUserLogs(c echo.Context) error {
	username := strings.ToLower(strings.TrimSpace(c.Param("username")))

	fromTime, toTime, err := parseFromTo(c.QueryParam("from"), c.QueryParam("to"), userHourLimit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userid := getUserid(username)

	var logResult chatLog

	orderBy := buildOrder(c)
	limit, err := buildLimit(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid limit")
	}

	selectFields := []string{"message", "dateOf(timeuuid)", "userid"}
	whereClauses := []string{"userid = ?", "timeuuid >= maxTimeuuid(?)", "timeuuid <= minTimeuuid(?)"}

	iter := cassandra.Query(
		buildQuery(selectFields, "logstv.user_messages", whereClauses, orderBy, limit),
		userid,
		fromTime,
		toTime,
	).Iter()

	var message chatMessage
	var ts time.Time
	var fetchedUserid int64
	var messageRaw string
	for iter.Scan(&messageRaw, &ts, &fetchedUserid) {
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
		return c.JSON(http.StatusInternalServerError, "Failure reading messages")
	}

	if c.Request().Header.Get("Content-Type") == "application/json" || c.QueryParam("type") == "json" {
		return writeJSONResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}
