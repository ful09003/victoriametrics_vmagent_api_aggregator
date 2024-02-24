package main

import "time"

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

func main() {
	panic("ope")
}
