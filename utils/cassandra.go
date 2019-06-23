package utils

import (
	"fmt"
	"net/http"
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
	qryString := "INSERT INTO data.info (payload, city_name, timestamp) VALUES (?, ?, ?) USING TTL 21600;"
	CassandraSession.Query(qryString, string(bytes), city_name, time.Now()).Exec()
}
func findOrCreate(cityName string) int {
	existQryString := "select count(id) as exists from cron_mt where city_name=?"
	result := map[string]interface{}{}
	CassandraSession.Query(existQryString, cityName).MapScan(result)
	if result["exists"].(int64) == 0 {
		insertQryString := "INSERT INTO data.cron_mt (city_name, id) VALUES (?,?)"
		err := CassandraSession.Query(insertQryString, cityName, gocql.TimeUUID()).Exec()
		if err != nil {
			fmt.Printf("FindOrCreate: %v\n", err)
			return http.StatusUnprocessableEntity
		}
		return http.StatusCreated
	} else {
		return http.StatusOK
	}
}
