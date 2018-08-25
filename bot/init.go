package main

import (
	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/gempir/logstv/common"
)

var queries = []string{
	`CREATE  KEYSPACE IF NOT EXISTS streamlogs
	WITH REPLICATION = { 
		'class' : 'SimpleStrategy', 
		'replication_factor' : 1 
	};`,
	`CREATE TABLE IF NOT EXISTS streamlogs.messages (
		channelId bigint,
		userId bigint,
		message text,
		timestamp timestamp,
		PRIMARY KEY (channelId, userId, timestamp)
	);`,
	`CREATE TABLE IF NOT EXISTS streamlogs.channels (
		userId bigint,
		username text,
		PRIMARY KEY (userId, username)
	);`,
	`CREATE INDEX IF NOT EXISTS channels_username_index
	ON streamlogs.channels (username)`,
}

func startup() {
	tClient = twitch.NewClient("justinfan123123", "oauth:123123123")

	hosts := strings.Split(common.GetEnv("DBHOSTS"), ",")
	cassandra, err := common.NewDatabaseSession(hosts)
	if err != nil {
		panic(err)
	}

	for _, query := range queries {
		err = cassandra.Query(query).Exec()
		if err != nil {
			panic(err)
		}
	}
}
