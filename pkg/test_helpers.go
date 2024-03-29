package pkg

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

const (
	happySampleText = `{"status":"success","data":{"activeTargets":[{"discoveredLabels":{"__address__":"nodeexporter:9100","__metrics_path__":"/metrics","__scheme__":"http","__scrape_interval__":"5s","__scrape_timeout__":"5s","job":"int_node_exporter"},"labels":{"instance":"nodeexporter:9100","job":"int_node_exporter"},"scrapePool":"int_node_exporter","scrapeUrl":"http://nodeexporter:9100/metrics","lastError":"","lastScrape":"2024-02-24T02:43:09.715Z","lastScrapeDuration":0.143,"lastSamplesScraped":918,"health":"up"}],"droppedTargets":[]}}`
	happyScrapeTime = `2024-02-24T02:43:09.715Z`
)

func happyResponse() VMAgentAPIResponse {
	t, _ := time.Parse(time.RFC3339, happyScrapeTime)
	return VMAgentAPIResponse{
		Data: Data{
			ActiveTargets: []VMAgentAPITarget{
				{
					DiscoveredLabels: map[string]string{
						"__address__":         "nodeexporter:9100",
						"__metrics_path__":    "/metrics",
						"__scheme__":          "http",
						"__scrape_interval__": "5s",
						"__scrape_timeout__":  "5s",
						"job":                 "int_node_exporter",
					},
					Labels: map[string]string{
						"job":      "int_node_exporter",
						"instance": "nodeexporter:9100",
					},
					ScrapePool:         "int_node_exporter",
					ScrapeURL:          "http://nodeexporter:9100/metrics",
					LastError:          "",
					LastScrape:         t,
					LastScrapeDuration: 0.143,
					LastSamplesScraped: 918,
					Health:             "up",
				},
			},
			DroppedTargets: []VMAgentAPITarget{},
		},
		Status: "success",
	}
}

func genTestHappyPathServerReq(t *testing.T) (*httptest.Server, *http.Request) {
	t.Helper()
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(happySampleText))
	}

	http.HandleFunc("/happy", handler)
	srv := httptest.NewServer(http.DefaultServeMux)
	u, err := url.Parse(srv.URL + "/happy")
	assert.NilError(t, err)

	return srv, &http.Request{URL: u}
}
