package api

import (
	"fmt"
	"net/http"
	"ttv-statistics/handlers"
)

const (
	apiName            = "ttv-statistics"
	getVideoStatistics = "getstreamervideostatistics"
)

var (
	EndpointMapping = map[string]func(w http.ResponseWriter, r *http.Request){
		fmt.Sprintf("/%s/%s/{%s}", apiName, getVideoStatistics, handlers.UserNamePathParam): handlers.GetStreamerVideoStatistics,
	}
)
