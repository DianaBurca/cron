package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
)

type CassandraConfig struct {
	Hosts                        []string
	Port                         int
	Keyspace, Username, Password string
}

var CassandraSession *gocql.Session
var CassandraCfg CassandraConfig

//EstablishConnection ...
func EstablishConnection() {
	var host string
	var password string
	if host = os.Getenv("host"); host == "" {
		panic("could not extract host")
	}
	if password = os.Getenv("pass"); password == "" {
		panic("could not extract password")
	}

	CassandraCfg.Hosts = append(CassandraCfg.Hosts, host)
	cluster := gocql.NewCluster()
	cluster.Hosts = CassandraCfg.Hosts
	cluster.Port = 9042
	cluster.Keyspace = "data"
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		CassandraSession = nil
		fmt.Println("could not connect to casssandra:", err)
	} else {
		CassandraSession = session
	}
}

// StoreInfo ...
func StoreInfo(bytes []byte, city_name string) {
	qryString := "INSERT INTO data.info (payload, city_name, timestamp) VALUES (?, ?, ?)"
	CassandraSession.Query(qryString, string(bytes), city_name, time.Now()).Exec()
}
