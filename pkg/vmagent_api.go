package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
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
		return vmRes, err
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

type VMAgentAPICollection struct {
	m    *sync.Mutex
	c    map[string]*VMAgentAPICollector
	data map[string]VMAgentAPIResponse
}

// NewVMAgentCollection initializes a new VMAgentAPICollection, used to hold data from multiple vmagent APIs
func NewVMAgentCollection(discovery VMAgentDiscoverer) (*VMAgentAPICollection, error) {
	c := &VMAgentAPICollection{
		m:    &sync.Mutex{},
		c:    map[string]*VMAgentAPICollector{},
		data: make(map[string]VMAgentAPIResponse),
	}
	endpoints, err := discovery.DiscoverEndpoints()
	if err != nil {
		return nil, err
	}
	for _, e := range endpoints {
		collector, err := NewVMAgentAPICollector(e, http.DefaultClient)
		if err != nil {
			return c, err
		}
		c.c[e] = collector
	}

	return c, nil
}

type VMAgentAPICollectionError struct {
	endpoint string
	err      error
}

func (e VMAgentAPICollectionError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.endpoint, e.err)
}

func (e VMAgentAPICollectionError) Unwrap() error {
	return e.err
}

func (v *VMAgentAPICollection) CollectAll() []error {
	errs := make([]error, 0)
	for endpoint, collector := range v.c {
		d, err := collector.Collect()
		v.m.Lock()
		if err != nil {
			v.data[endpoint] = VMAgentAPIResponse{}
			errs = append(errs, VMAgentAPICollectionError{
				endpoint: endpoint,
				err:      err,
			})
		} else {
			v.data[endpoint] = d
		}
		v.m.Unlock()
	}

	return errs
}

func (v *VMAgentAPICollection) Data() map[string]VMAgentAPIResponse {
	return v.data
}

// Reconcile accepts a list of new endpoints a VMAgentAPICollection should track.
// New entries are added to the collection, while endpoints in the VMAgentAPICollection but not in the new list are removed
func (v *VMAgentAPICollection) Reconcile(newEndpoints []string) error {
	v.m.Lock()
	defer v.m.Unlock()
	newCollectors := map[string]struct{}{}

	for _, e := range newEndpoints {
		newCollectors[e] = struct{}{}
	}

	for e, _ := range v.c {
		// If existing collectors map contains an endpoint not in new collectors, remove it
		if _, exists := newCollectors[e]; !exists {
			delete(v.c, e)
			delete(v.data, e)
		}
	}
	for e, _ := range newCollectors {
		// If new collection contains an entry not in existing, add it
		if _, exists := v.c[e]; !exists {
			newCollector, err := NewVMAgentAPICollector(e, http.DefaultClient)
			if err != nil {
				return err
			}
			v.c[e] = newCollector
			v.data[e] = VMAgentAPIResponse{}
		}
	}
	return nil
}
