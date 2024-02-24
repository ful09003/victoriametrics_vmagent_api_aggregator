package pkg

import (
	"net/http"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func Test_fetchVMAgentTargets(t *testing.T) {
	hst, _ := time.Parse(time.RFC3339, happyScrapeTime)
	type args struct {
		c *http.Client
	}
	tests := []struct {
		name    string
		args    args
		want    VMAgentAPIResponse
		wantErr error
	}{
		{
			name: "happy path",
			args: args{
				c: http.DefaultClient,
			},
			want: VMAgentAPIResponse{
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
							LastScrape:         hst,
							LastScrapeDuration: 0.143,
							LastSamplesScraped: 918,
							Health:             "up",
						},
					},
					DroppedTargets: []VMAgentAPITarget{},
				},
				Status: "success",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, req := genTestHappyPathServerReq(t)
			defer srv.Close()

			got, err := fetchVMAgentTargets(tt.args.c, req)
			assert.Equal(t, err, tt.wantErr)
			assert.DeepEqual(t, got, tt.want)
		})
	}
}
