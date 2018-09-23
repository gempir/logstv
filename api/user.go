package main

import (
	"github.com/gocql/gocql"
)

var userCache = make(map[string]int64)

func getUsernameByUserid(userid int64) string {
	var username string

	cassandra.Query(`SELECT username FROM logstv.users WHERE userid = ? LIMIT 1`, userid).Consistency(gocql.One).Scan(&username)

	return username
}

func getUseridbyUsername(username string) int64 {
	var userid int64

	cassandra.Query(`SELECT userid FROM logstv.users WHERE username = ? LIMIT 1`, username).Consistency(gocql.One).Scan(&userid)

	return userid
}

func getUserid(username string) int64 {
	userid, ok := userCache[username]
	if !ok {
		userid = getUseridbyUsername(username)
		userCache[username] = userid
	}

	return userid
}
