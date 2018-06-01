package rrdtool

import (
	"log"
	"math"
	"time"
	"sync"

	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/rrdlite"
	"github.com/toolkits/file"

	"github.com/shenxingwuying/open-falcon-study/graph/g"
	"github.com/shenxingwuying/open-falcon-study/graph/store"

	"github.com/influxdata/influxdb/client/v2"
)

func WriteInfluxdb(filename string, items []*cmodel.GraphItem) error {
	cfg := g.Config()

	// temp, because of type cast
	// Create a new HTTPClient
	// @begin
	cfg := g.Config()
	influxdbClient , err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.66.0.220:8086",
		Username: cfg.Username,
		Password: cfg.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer influxdbClient.Close()
	// @end

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  cfg.Database,
		Precision: "s",
	})
	if err != nil {
		log.Fatal("new batch point err", err)
	}

	for _, item := range items {
		tags := map[string]string {
			"endpoint": item.Endpoint,
			"counter": item.Metric,
			"tags": item.Tags,
			"type": item.DsType,
			"step": item.Step,
		}
		fields := map[string] interface{} {
			"value": item.Value,
		}
		pt, err := client.NewPoint("open_falcon_table", tags, fields, item.Timestamp)
		if err != nil {
			log.Fatal("new point error", err)
		}
		bp.AddPoint(pt)
	}

	if err := influxdbClient(client.client).Write(bp); err != nil {
		log.Fatal(err)
	}

	return nil
}

//func getInfluxdbClient() (client *client.Client) {
//	rwLock.Lock()
//	defer rwLock.Unlock()
//
//	if (client == nil) {
//		client, err := initDB()
//		if (err != nil) {
//			log.Fatal("get influxdb client error")
//		}
//	}
//	return client
//}
//
//func closeClient() {
//	rwLock.Lock()
//	defer rwLock.Unlock()
//
//	if influxdbClient != nil {
//		client.Close()
//	}
//}