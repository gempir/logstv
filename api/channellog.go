package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var channelHourLimit = 24.0

func getChannelLogs(c echo.Context) error {
	channel := strings.TrimSpace(strings.ToLower(c.Param("channel")))

	from := c.QueryParam("from")
	to := c.QueryParam("to")

	var fromTime time.Time
	var toTime time.Time

	if from == "" && to == "" {
		fromTime = time.Now().AddDate(0, -1, 0)
		toTime = time.Now()
	} else if from == "" && to != "" {
		var err error
		toTime, err = parseTimestamp(to)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Can't parse to timestamp: %s", err))
		}
		fromTime = toTime.AddDate(0, -1, 0)
	} else if from != "" && to == "" {
		var err error
		fromTime, err = parseTimestamp(from)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Can't parse from timestamp: %s", err))
		}
		toTime = fromTime.AddDate(0, 1, 0)
	} else {
		var err error

		fromTime, err = parseTimestamp(from)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Can't parse from timestamp: %s", err))
		}
		toTime, err = parseTimestamp(to)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Can't parse to timestamp: %s", err))
		}

		if toTime.Sub(fromTime).Hours() > channelHourLimit {
			return c.String(http.StatusBadRequest, "Timespan too big")
		}
	}

	channelid := getUserid(channel)

	var logResult chatLog
	var iter *gocql.Iter

	_, reverse := c.QueryParams()["reverse"]
	if reverse {
		limit := c.QueryParam("limit")
		if limit != "" {
			limitInt, err := strconv.Atoi(limit)
			if err != nil || limitInt < 1 {
				return c.JSON(http.StatusBadRequest, "Invalid limit")
			}

			iter = cassandra.Query(`
			SELECT message, timestamp, userid
			FROM logstv.channel_messages 
			WHERE channelid = ? 
			AND timestamp >= ? 
			AND timestamp <= ?
			ORDER BY timestamp DESC
			LIMIT ?`,
				channelid,
				fromTime,
				toTime,
				limitInt).Iter()
		} else {
			iter = cassandra.Query(`
			SELECT message, timestamp, userid
			FROM logstv.channel_messages 
			WHERE channelid = ? 
			AND timestamp >= ? 
			AND timestamp <= ?
			ORDER BY timestamp DESC`,
				channelid,
				fromTime,
				toTime).Iter()
		}
	} else {
		limit := c.QueryParam("limit")
		if limit != "" {
			limitInt, err := strconv.Atoi(limit)
			if err != nil || limitInt < 1 {
				return c.JSON(http.StatusBadRequest, "Invalid limit")
			}

			iter = cassandra.Query(`
			SELECT message, timestamp, userid
			FROM logstv.channel_messages 
			WHERE channelid = ? 
			AND timestamp >= ? 
			AND timestamp <= ?
			ORDER BY timestamp ASC
			LIMIT ?`,
				channelid,
				fromTime,
				toTime,
				limitInt).Iter()
		} else {
			iter = cassandra.Query(`
			SELECT message, timestamp, userid
			FROM logstv.channel_messages 
			WHERE channelid = ? 
			AND timestamp >= ? 
			AND timestamp <= ?
			ORDER BY timestamp ASC`,
				channelid,
				fromTime,
				toTime).Iter()
		}
	}

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
		return c.JSON(http.StatusOK, logResult)
	}

	return c.String(http.StatusOK, buildTextChatLog(logResult))
}
