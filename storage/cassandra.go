package storage

import (
	"github.com/gocql/gocql"
	"time"
)

type CassandraClient struct {
	Session *gocql.Session
}

func NewCassandraClient(hosts []string, keyspace string) (*CassandraClient, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS location_history (
			user_id text,
            timestamp timestamp,
            latitude double,
            longitude double,
            PRIMARY KEY (user_id, timestamp)
		)
	`).Exec()
	if err != nil {
		return nil, err
	}

	return &CassandraClient{Session: session}, nil
}

func (c *CassandraClient) SaveLocation(userId string, lat, lon float64, timestamp int64) error {
	return c.Session.Query(
		"INSERT INTO location_history (user_id, timestamp, latitude, longitude) VALUES(?, ?, ?, ?)",
		userId, time.Unix(timestamp, 0), lat, lon,
	).Exec()
}
