package rrdtool

import (
    "fmt"
	"log"
    "time"
    "bytes"
    "strconv"

	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"
	"github.com/open-falcon/graph/g"
	"github.com/influxdata/influxdb/client/v2"
)

var influxdbMasterClient client.Client
var influxdbSlaveClient client.Client

func InitInfluxdbClient() {
	cfg := g.Config()

	// temp, because of type cast
	// Create a new HTTPClient
	// @begin
    influxdbMasterClient, err1 := client.NewHTTPClient(client.HTTPConfig{
		Address:  cfg.Influxdb.Address.Master,
		Username: cfg.Influxdb.Username,
		Password: cfg.Influxdb.Password,
	})

	influxdbSlaveClient, err2 := client.NewHTTPClient(client.HTTPConfig{
		Address:  cfg.Influxdb.Address.Slave,
		Username: cfg.Influxdb.Username,
		Password: cfg.Influxdb.Password,
	})

	if err1 != nil && err2 != nil {
		log.Fatal(err1)
	}
}

func WriteInfluxdb(filename string, items []*cmodel.GraphItem, isMaster bool) error {
	cfg := g.Config()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  cfg.Influxdb.Database,
		Precision: "s",
	})
	if err != nil {
		log.Fatal("new batch point err", err)
	}

	for _, item := range items {
		// counter := item.Metric+"/"+sort(item.Tags)
		counter := item.Metric
		if len(item.Tags) > 0 {
			counter = fmt.Sprintf("%s/%s", counter, cutils.SortedTags(item.Tags))
		}
		tags := map[string]string{
			"endpoint": item.Endpoint,
			"counter":  counter,
		}
		fields := map[string]interface{}{
			"value": item.Value,
		}
		pt, err := client.NewPoint(cfg.Influxdb.Tablename, tags, fields, time.Unix(item.Timestamp, 0))
		if err != nil {
			log.Fatal("new point error", err)
		}
		bp.AddPoint(pt)
	}

	if isMaster {
		if err := influxdbMasterClient.Write(bp); err != nil {
			log.Println(err)
		}
	} else {
		if err := influxdbSlaveClient.Write(bp); err != nil {
			log.Println(err)
		}
	}


	return nil
}


func ReadInfluxdb(endpoint string, counter string, cf string, start int64, end int64, step int, isMaster bool) ([]client.Result, error) {
	cfg := g.Config()

	var buffer bytes.Buffer
	buffer.WriteString("select time, value from open_falcon_table where endpoint='"+endpoint)
	buffer.WriteString("' and counter= '" + counter)
	buffer.WriteString("' and time >= " + strconv.FormatInt(1000000000 * start, 10))
	buffer.WriteString(" and time <= " + strconv.FormatInt(1000000000 * end, 10))

	querySql := buffer.String()
	q := client.Query {
		Command: querySql,
		Database: cfg.Influxdb.Database,
	}

	var res []client.Result
	if isMaster {
		if response, err := influxdbMasterClient.Query(q); err == nil {
			if response.Error() != nil {
				log.Println(response.Error())
			}
			res = response.Results
		} else {
			return nil, nil
		}
	} else {
		if response, err := influxdbSlaveClient.Query(q); err == nil {
			if response.Error() != nil {
				log.Println(response.Error())
			}
			res = response.Results
		} else {
			return nil, nil
		}
	}

	return res, nil
}