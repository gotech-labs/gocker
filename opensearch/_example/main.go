package main

import (
	"net/http"

	"github.com/gotech-labs/core/log"
	"github.com/gotech-labs/gocker/opensearch"
	"github.com/gotech-labs/gocker/opensearch/dashboard"
)

func main() {
	var (
		tag = "1.3.1"
	)
	aos := opensearch.New(tag,
		"analysis-kuromoji",
		"analysis-icu",
	)
	callAPI(aos.RestAPIEndpoint())
	println(aos.AcceptDashboardURL())

	dashboard := dashboard.New(tag, aos.AcceptDashboardURL())
	callAPI(dashboard.RestAPIEndpoint())

	log.Info().
		Str("opensearch", aos.RestAPIEndpoint()).
		Str("dashboard", dashboard.RestAPIEndpoint()).
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
