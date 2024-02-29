package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ful09003/victoriametrics_vmagent_api_aggregator/pkg"
)

var (
	flagVMAgentTargets       = flag.String("targetdiscovery", "", "file path containing targets or name of env var to reference")
	flagTargetsWatchInterval = flag.Duration("discoveryinterval", 10*time.Second, "duration to update discovered targets list")
)

func main() {
	flag.Parse()
	disco := newDiscovery(*flagVMAgentTargets)

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
