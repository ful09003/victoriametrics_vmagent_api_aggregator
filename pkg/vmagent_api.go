package pkg

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type VMAgentAPIResponse struct {
	Data   Data   `json:"data"`
	Status string `json:"status"`
}
type VMAgentAPITarget struct {
	DiscoveredLabels   map[string]string `json:"discoveredLabels"`
	Labels             map[string]string `json:"labels"`
	ScrapePool         string            `json:"scrapePool"`
	ScrapeURL          string            `json:"scrapeUrl"`
	LastError          string            `json:"lastError"`
	LastScrape         time.Time         `json:"lastScrape"`
	LastScrapeDuration float64           `json:"lastScrapeDuration"`
	LastSamplesScraped int               `json:"lastSamplesScraped"`
	Health             string            `json:"health"`
}
type Data struct {
	ActiveTargets  []VMAgentAPITarget `json:"activeTargets"`
	DroppedTargets []VMAgentAPITarget `json:"droppedTargets"`
}

func fetchVMAgentTargets(c *http.Client, r *http.Request) (VMAgentAPIResponse, error) {
	vmRes := VMAgentAPIResponse{}

	resp, err := c.Do(r)
	if err != nil {
		return vmRes, nil
	}
	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return vmRes, nil
	}

	if err := json.Unmarshal(resBytes, &vmRes); err != nil {
		return VMAgentAPIResponse{}, err
	}

	return vmRes, nil
}

type VMAgentAPICollector struct {
	origEndpoint string
	c            *http.Client
}

// NewVMAgentAPICollector returns a collector for a vmagent API endpoint. It is expected that the full URL to
// the vmagent targets API is provided
func NewVMAgentAPICollector(s string, c *http.Client) (*VMAgentAPICollector, error) {
	return &VMAgentAPICollector{
		origEndpoint: s,
		c:            c,
	}, nil
}

// Collect returns a vmagent API response
func (v *VMAgentAPICollector) Collect() (VMAgentAPIResponse, error) {
	u, err := url.Parse(v.origEndpoint)
	if err != nil {
		return VMAgentAPIResponse{}, err
	}
	req := &http.Request{URL: u}
	return fetchVMAgentTargets(v.c, req)
}
