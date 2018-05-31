package g

import (
	"error"
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
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.66.0.220:8086",
		Username: cfg.username,
		Password: cfg.password,
	})
	if err != nil {
		log.Fatal(err)
	}
	client = c
	return c, err
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
