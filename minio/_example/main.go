package main

import (
	"net/http"

	"github.com/gotech-labs/core/log"
	"github.com/gotech-labs/gocker/minio"
)

func main() {
	s3 := minio.New("latest", "test")
	callAPI(s3.RestAPIEndpoint() + "/minio/health/live")

	log.Info().
		Str("api endpoint", s3.RestAPIEndpoint()).
		Str("console address", s3.ConsoleAddress()).
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
