package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var hourLimit = 744.0

type message struct {
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type userlog struct {
	Username string    `json:"username"`
	Messages []message `json:"message"`
}

func getUserLogs(c echo.Context) error {
	channel := strings.TrimSpace(strings.ToLower(c.Param("channel")))
	username := strings.ToLower(strings.TrimSpace(c.Param("username")))

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

		if toTime.Sub(fromTime).Hours() > hourLimit {
			return c.String(http.StatusBadRequest, "Timespan too big")
		}
	}

	channelid := getUserid(channel)
	userid := getUserid(username)

	var userlogResult userlog
	userlogResult.Username = username

	var message message
	iter := cassandra.Query(`
	 SELECT message, timestamp
	 FROM logstv.messages 
	 WHERE userid = ? 
	 AND channelid = ? 
	 AND timestamp >= ? 
	 AND timestamp <= ?`,
		userid,
		channelid,
		fromTime,
		toTime).Iter()
	for iter.Scan(&message.Text, &message.Timestamp) {

		userlogResult.Messages = append(userlogResult.Messages, message)
	}
	if err := iter.Close(); err != nil {
		log.Error(err)
	}

	if c.Request().Header.Get("Content-Type") == "application/json" {
		return c.JSON(http.StatusOK, userlogResult)
	}

	return c.String(http.StatusOK, buildTextUserlog(userlogResult))
}

func buildTextUserlog(ulog userlog) string {
	var text string

	for _, message := range ulog.Messages {
		text += fmt.Sprintf("[%s] %s: %s\r\n", message.Timestamp.Format("2006-01-2 15:04:05 UTC"), ulog.Username, message.Text)
	}

	return text
}

func getUserid(username string) int64 {
	userid, ok := userCache[username]
	if !ok {
		userid = getUseridbyUsername(username)
		userCache[username] = userid
	}

	return userid
}

func parseTimestamp(timestamp string) (time.Time, error) {

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil
}
