package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var userHourLimit = 744.0

type userlog struct {
	Username string        `json:"username"`
	Messages []userMessage `json:"message"`
}

type userMessage struct {
	Text      string    `json:"text"`
	Timestamp timestamp `json:"timestamp"`
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

		if toTime.Sub(fromTime).Hours() > userHourLimit {
			return c.String(http.StatusBadRequest, "Timespan too big")
		}
	}

	channelid := getUserid(channel)
	userid := getUserid(username)

	var userlogResult userlog
	userlogResult.Username = username

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

	var message userMessage
	var ts time.Time
	for iter.Scan(&message.Text, &ts) {
		message.Timestamp = timestamp{ts}

		userlogResult.Messages = append(userlogResult.Messages, message)
	}
	if err := iter.Close(); err != nil {
		log.Error(err)
	}

	_, reverse := c.QueryParams()["reverse"]
	if reverse {
		sort.SliceStable(userlogResult.Messages, func(i, j int) bool {
			return userlogResult.Messages[i].Timestamp.Unix() > userlogResult.Messages[j].Timestamp.Unix()
		})
	}

	limit := c.QueryParam("limit")
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid limit")
		}
		userlogResult.Messages = userlogResult.Messages[:limitInt]
	}

	if c.Request().Header.Get("Content-Type") == "application/json" || c.QueryParam("type") == "json" {
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

func parseTimestamp(timestamp string) (time.Time, error) {

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(i, 0), nil
}
