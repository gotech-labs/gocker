package main

import (
	"net/http"

	"github.com/gotech-labs/core/log"
	"github.com/gotech-labs/gocker/elasticsearch"
	"github.com/gotech-labs/gocker/elasticsearch/kibana"
)

func main() {
	var (
		tag = "8.2.0"
	)
	es := elasticsearch.New(tag,
		"analysis-kuromoji",
		"analysis-icu",
	)
	callAPI(es.RestAPIEndpoint())

	dashboard := kibana.New(tag, es.AcceptDashboardURL())
	callAPI(dashboard.RestAPIEndpoint())

	log.Info().
		Str("elasticsearch", es.RestAPIEndpoint()).
		Str("kibana", dashboard.RestAPIEndpoint()).
		Send()
}

func callAPI(endpoint string) {
	log.Info().Msg(endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect server")
	}
	if resp.StatusCode != 200 {
		log.Panic().Msgf("Failed to access server [status=%d]", resp.StatusCode)
	}
}
