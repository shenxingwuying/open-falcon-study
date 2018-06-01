package rrdtool

import (
	"log"
    "time"

	cmodel "github.com/open-falcon/common/model"

	"github.com/shenxingwuying/open-falcon-study/graph/g"
	"github.com/influxdata/influxdb/client/v2"
)

func WriteInfluxdb(filename string, items []*cmodel.GraphItem) error {
	cfg := g.Config()

	// temp, because of type cast
	// Create a new HTTPClient
	// @begin
	influxdbClient , err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.66.0.220:8086",
		Username: cfg.Influxdb.Username,
		Password: cfg.Influxdb.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer influxdbClient.Close()
	// @end

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  cfg.Influxdb.Database,
		Precision: "s",
	})
	if err != nil {
		log.Fatal("new batch point err", err)
	}

	for _, item := range items {
        // counter := item.Metric+"/"+sort(item.Tags)
        counter := item.Metric+"/"
		tags := map[string]string {
			"endpoint": item.Endpoint,
			"counter": counter,
		}
		fields := map[string] interface{} {
			"value": item.Value,
			"type": item.DsType,
			"step": item.Step,
		}
		pt, err := client.NewPoint("open_falcon_table", tags, fields, time.Unix(item.Timestamp, 0))
		if err != nil {
			log.Fatal("new point error", err)
		}
		bp.AddPoint(pt)
	}

	if err := influxdbClient.Write(bp); err != nil {
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
