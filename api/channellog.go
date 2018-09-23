package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var channelHourLimit = 24.0

type channelLog struct {
	Username string           `json:"username"`
	Messages []channelMessage `json:"message"`
}

type channelMessage struct {
	Text      string    `json:"text"`
	Username  string    `json:"username"`
	Timestamp timestamp `json:"timestamp"`
}

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

	var channelLogResult channelLog

	iter := cassandra.Query(`
	 SELECT message, timestamp, userid
	 FROM logstv.messages
	 WHERE channelid = ?
	 AND timestamp >= ?
	 AND timestamp <= ?`,
		channelid,
		fromTime,
		toTime).Iter()

	var message channelMessage
	var ts time.Time
	var userid int64
	for iter.Scan(&message.Text, &ts, &userid) {
		message.Timestamp = timestamp{ts}
		message.Username = getUsernameByUserid(userid)

		channelLogResult.Messages = append(channelLogResult.Messages, message)
	}
	if err := iter.Close(); err != nil {
		log.Error(err)
	}

	if c.Request().Header.Get("Content-Type") == "application/json" || c.QueryParam("type") == "json" {
		return c.JSON(http.StatusOK, channelLogResult)
	}

	return c.String(http.StatusOK, "buildTextUserlog(channelLogResult)")
}

// func buildTextUserlog(ulog userlog) string {
// 	var text string

// 	for _, message := range ulog.Messages {
// 		text += fmt.Sprintf("[%s] %s: %s\r\n", message.Timestamp.Format("2006-01-2 15:04:05 UTC"), ulog.Username, message.Text)
// 	}

// 	return text
// }
