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

func getCurrentUserLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)
	year := time.Now().Year()
	month := time.Now().Month().String()
	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))

	redirectURL := fmt.Sprintf("/channel/%s/user/%s/%d/%s", channel, username, year, month)
	return c.Redirect(303, redirectURL)
}

func getDatedUserLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)
	year := c.Param("year")
	month := strings.Title(c.Param("month"))
	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))

	if year == "" || month == "" {
		year = strconv.Itoa(time.Now().Year())
		month = time.Now().Month().String()
	}

	userid, ok := userCache[username]
	if !ok {
		userid = getUseridbyUsername(username)
	}

	channelid, ok := userCache[channel]
	if !ok {
		channelid = getUseridbyUsername(channel)
	}

	userCache[username] = userid
	userCache[channel] = channelid

	var logResult string

	var message string
	iter := cassandra.Query("SELECT message FROM streamlogs.messages WHERE userid = ? AND channelid = ?", userid, channelid).Iter()
	for iter.Scan(&message) {
		logResult += fmt.Sprintf("%s: %s\r\n", username, message)
	}
	if err := iter.Close(); err != nil {
		log.Error(err)
	}

	return c.String(http.StatusOK, logResult)
}
