package main

import (
	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/gempir/logstv/common"
)

var queries = []string{
	`CREATE  KEYSPACE IF NOT EXISTS logstv
	WITH REPLICATION = { 
		'class' : 'SimpleStrategy', 
		'replication_factor' : 1 
	};`,
	`CREATE TABLE IF NOT EXISTS logstv.messages (
		timestamp timestamp,
		channelid bigint,
		userid bigint,
		message text,
		PRIMARY KEY ((channelid, userid), timestamp)
	) WITH CLUSTERING ORDER BY (timestamp DESC);`,
	`CREATE TABLE IF NOT EXISTS logstv.channel_messages (
		timestamp timestamp,
		channelid bigint,
		userid bigint,
		message text,
		PRIMARY KEY ((channelid), timestamp)
	) WITH CLUSTERING ORDER BY (timestamp DESC);`,
	`CREATE TABLE IF NOT EXISTS logstv.user_messages (
		timestamp timestamp,
		channelid bigint,
		userid bigint,
		message text,
		PRIMARY KEY ((userid), timestamp)
	) WITH CLUSTERING ORDER BY (timestamp DESC);`,
	`CREATE TABLE IF NOT EXISTS logstv.channels (
		userid bigint,
		username text,
		PRIMARY KEY (userid, username)
	);`,
	`CREATE INDEX IF NOT EXISTS channels_username_index ON logstv.channels (username)`,
	`CREATE TABLE IF NOT EXISTS logstv.users (
		userid bigint,
		username text,
		PRIMARY KEY (userid, username)
	);`,
	`CREATE INDEX IF NOT EXISTS users_username_index ON logstv.users (username)`,
}

func startup() {
	tClient = twitch.NewClient("justinfan123123", "oauth:123123123")
	tClient.TLS = false
	tClient.IrcAddress = "127.0.0.1:3333"

	hosts := strings.Split(common.GetEnv("DBHOSTS"), ",")
	var err error
	cassandra, err = common.NewDatabaseSession(hosts)
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
