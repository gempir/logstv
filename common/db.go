package common

import "github.com/gocql/gocql"

func NewDatabaseSession(DBHosts []string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(DBHosts...)

	cassandra, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return cassandra, nil
}
