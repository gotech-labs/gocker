package opensearch_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gotech-labs/gocker/opensearch"
)

func TestLogger(t *testing.T) {
	container := New("1.3.1")
	defer container.Purge()

	endpoint := "http://localhost:" + container.Port("9200/tcp")
	resp, err := http.Get(endpoint)
	if assert.NoError(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
	}
}
