package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var userHourLimit = 744.0

func getUserLogs(c echo.Context) error {
	channel := strings.TrimSpace(strings.ToLower(c.Param("channel")))
	username := strings.ToLower(strings.TrimSpace(c.Param("username")))

	fromTime, toTime, err := parseFromTo(c.QueryParam("from"), c.QueryParam("to"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	channelid := getUserid(channel)
	userid := getUserid(username)

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
	whereClauses := []string{"userid = ?", "channelid = ?", "timestamp >= ?", "timestamp <= ?"}

	iter := cassandra.Query(
		buildQuery(selectFields, "logstv.messages", whereClauses, orderBy, limitInt),
		userid,
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
		return c.JSON(http.StatusInternalServerError, "Failure reading messages")
	}

	if c.Request().Header.Get("Content-Type") == "application/json" || c.QueryParam("type") == "json" {
		return writeJSONResponse(c, &logResult)
	}

	return writeTextResponse(c, &logResult)
}

func parseTimestamp(timestamp string) (time.Time, error) {

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil
}
