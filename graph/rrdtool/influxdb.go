package rrdtool

import (
	"error"
	"log"
	"math"
	"time"

	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/rrdlite"
	"github.com/toolkits/file"

	"github.com/shenxingwuying/open-falcon-study/graph/g"
	"github.com/shenxingwuying/open-falcon-study/graph/store"

	"github.com/influxdata/influxdb/client/v2"
)

func write_influxdb(filename string, items[] *cmodel.GraphItem) error {
	return nil;
}
