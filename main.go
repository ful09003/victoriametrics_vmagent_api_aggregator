package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ful09003/victoriametrics_vmagent_api_aggregator/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	flagVMAgentTargets       = flag.String("targetdiscovery", "", "file path containing targets or name of env var to reference")
	flagTargetsWatchInterval = flag.Duration("discoveryinterval", 10*time.Second, "duration to update discovered targets list")
)

var (
	lastScrapedVec *prometheus.GaugeVec
)

func main() {
	flag.Parse()
	disco := newDiscovery(*flagVMAgentTargets)

	lastScrapedVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "vma",
		Name:        "last_samples_scraped",
		Help:        "Total number of samples last collected per-job for each tracked vmagent instance",
		ConstLabels: nil,
	},
		[]string{"vmagent_instance", "job"})

	reg := prometheus.NewRegistry()
	reg.MustRegister(lastScrapedVec, collectors.NewGoCollector())
	pH := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	http.Handle("/debug/metrics", pH)
	go func() {
		if err := http.ListenAndServe(":18429", nil); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	collection, err := pkg.NewVMAgentCollection(disco)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	t := time.NewTicker(*flagTargetsWatchInterval)
	for range t.C {
		if err := updateCollection(collection, *flagVMAgentTargets); err != nil {
			log.Println(err)
		}
		if errs := collection.CollectAll(); len(errs) != 0 {
			log.Println(errs)
		}
		if err := updatemetrics(lastScrapedVec, collection); err != nil {
			log.Println(err)
		}
	}
}

func newDiscovery(s string) pkg.VMAgentDiscoverer {
	if _, err := os.Stat(s); err == nil {
		return pkg.NewFileDiscovery(s)
	}
	return pkg.NewEnvDiscovery(s)
}

func updateCollection(c *pkg.VMAgentAPICollection, fp string) error {
	b, err := os.ReadFile(fp)
	if err != nil {
		return err
	}
	newTargets := strings.Split(strings.TrimSpace(string(b)), ",")
	return c.Reconcile(newTargets)
}

func updatemetrics(c *prometheus.GaugeVec, col *pkg.VMAgentAPICollection) error {
	d := col.Data()
	for agent, res := range d {
		for _, a := range res.Data.ActiveTargets {
			c.With(prometheus.Labels{
				"vmagent_instance": agent,
				"job":              a.DiscoveredLabels["job"],
			}).Set(float64(a.LastSamplesScraped))
		}
	}
	return nil
}
