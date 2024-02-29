# victoriametrics_vmagent_api_aggregator (VAA)

# What??
An application allowing for aggregating data from a fleet of [vmagents](https://docs.victoriametrics.com/vmagent/)

# Why??
vmagent exposes some very nice, desirable information in its `/api/v1/targets` endpoint .
Unfortunately, gathering that data is an exercise left to the deployer or (human) operator of vmagent.
This library aims to make this exercise easier to do!

# How??
This library is designed to follow a simple mental model:

- A Collection, which is similar to a library of vmagents, which are each described by and data gathered via...
- A Collector, which encapsulates HTTP client behaviors to request and convert the `/api/v1/targets` endpoints into a
  usable Golang struct. Each Collector is given a specific vmagent endpoint via...
- A VMAgentDiscoverer, an interface type that can feed a list of distinct targets back to the Collection

So basically,
- Set up a Discoverer (support for env var, file based, and static-in-go-code exists)
- Use that to create a Collection
- Call the collection.CollectAll() on a timer, etc.

I'm not truly happy with the state of this library! It's not the safest, but is bare-bones enough to put out into the
world and accept ridicule, praise, etc. over :)

# Demo time!
There is an example main.go in this repo. To see a truly basic example of what VAA unlocks, do the following:
- in terminal 1: `cd integation && podman-compose up`
- in terminal 2: `go build && ./victoriametrics_vmagent_api_aggregator -discoveryinterval=1s -targetdiscovery=$(pwd)/integration/vmagent_targets`
- in browser, navigate to `localhost:18429/debug/metrics`, looking for a `vma_last_samples_scraped` series

Gathering this data otherwise (total samples scraped by job for each vmagent in your fleet) would be otherwise 
highly manual today, and you would need to cURL, etc. the individual endpoints (localhost:8429, localhost:8430, etc.)

See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/4018 for one potential future
application of this library (building your own GUI which can display a true global aggregated view of your vmagent cluster)

