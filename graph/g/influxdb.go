package g

import (
	"log"
	"time"
	"sync"
	"github.com/influxdata/influxdb/client/v2"
)

var (
	rwLock sync.RWMutex
	client *client.Client
)

func initDB() (c *client.Client, err error) {
	// Create a new HTTPClient
	cfg := g.Config()
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.66.0.220:8086",
		Username: cfg.Username,
		Password: cfg.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	client = c
	return c, err
}

func WriteInfluxdb(filename string, items []*cmodel.GraphItem) error {
	cfg := g.Config()
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  cfg.Database,
		Precision: "s",
	})
	if err != nil {
		log.Fatal("new batch point err", err)
	}

	
}

func client() (client *client.Client) {
	rwLock.Lock()
	defer rwLock.Unlock()

	if (client == nil) {
		client, err := initDB()
		if (err != nil) {
			log.Fatal("get influxdb client error")
		}
	}
	return client
}
