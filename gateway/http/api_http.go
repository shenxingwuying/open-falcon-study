package http

import (
	"encoding/json"
	"net/http"

	cmodel "github.com/open-falcon/common/model"

	"github.com/open-falcon/gateway/g"
	trpc "github.com/open-falcon/gateway/receiver/rpc"
)

func configApiHttpRoutes() {
	http.HandleFunc("/api/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "blank body", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var metrics []*cmodel.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "decode error", http.StatusBadRequest)
			return
		}

		reply := &g.TransferResp{}
		trpc.RecvMetricValues(metrics, reply, "http")

		RenderDataJson(w, reply)
	})
}
