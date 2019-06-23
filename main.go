package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DianaBurca/cron/utils"
	"github.com/gin-gonic/gin"
)

func doRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func main() {
	utils.EstablishConnection()
	if utils.CassandraSession == nil {
		panic("could not establish session to cassandra")
	}
	defer utils.CassandraSession.Close()

	go func() {
		for {
			fmt.Println("will strart")
			infoLoader()
			time.Sleep(10 * time.Minute)
		}

	}()

	driver := gin.Default()

	driver.PUT("/store", utils.StoreHandler)
	driver.GET("/.well-known/live", utils.Health)
	driver.GET("/.well-known/ready", utils.Health)

	driver.Run()
}

func infoLoader() {
	var cities []string
	m := map[string]interface{}{}

	qryString := "SELECT city_name FROM data.cron_mt"

	iterable := utils.CassandraSession.Query(qryString).Iter()
	for iterable.MapScan(m) {
		cities = append(cities, m["city_name"].(string))
		m = map[string]interface{}{}
	}

	fmt.Println("Cities: ", cities)
	for _, city := range cities {
		resp, err := doRequest(fmt.Sprintf("http://fetcher/fetch?city=%s", city))
		if err == nil {
			fmt.Printf("Will store for %s", city)
			var bodyBytes []byte
			if resp.StatusCode == http.StatusOK {
				bodyBytes, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
				}
				utils.StoreInfo(bodyBytes, city)
			}
		} else {
			fmt.Println("err: ", err)
		}
	}
}
