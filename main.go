package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"

	log "github.com/sirupsen/logrus"
)

func init() {
	// log json because cloud
	if os.Getenv("ENVIRONMENT") == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
	log.SetOutput(os.Stdout)
}

func main() {
	var (
		project string
	)

	// cli flags
	flag.StringVar(&project, "project", "", "Google Cloud Project")
	flag.Parse()

	if project == "" {
		fmt.Println("Usage:\ngo-check-dns -project gcp-project")
		os.Exit(1)
	}

	// start datadog client
	c, err := statsd.New(os.Getenv("DATADOG_URL"))
	if err != nil {
		log.Fatalf("datadog could not start: %v", err)
	}
	log.Infof("Logging to: %s", os.Getenv("DATADOG_URL"))

	// add stat namespace
	c.Namespace = "go-check-dns."

	// get domains from Google Cloud DNS
	log.Infoln("Fetching domains...")
	domains := getZones(project)

	// iterate over the domains and do things
	log.Debugf("Performing WHOIS on: %s", domains)
	for _, domain := range domains {
		// get domain expiry
		log.Debugf("Looking up domain: %s", domain)
		expiry, err := lookup(domain)
		if err != nil {
			// skip this domain if we can't look it up
			log.Errorf("Could not get expiry for %s: %v", domain, err)
			break
		}

		// calculate time and debug
		timeDiff := expiry.Sub(time.Now())
		daysLeft := float64(timeDiff.Hours() / 24)
		log.WithFields(log.Fields{
			"domain":    domain,
			"expiry":    expiry,
			"days left": daysLeft,
		}).Debug("Domain Parsed!")

		// submit datadog metric
		var tags []string

		// idk if a tag is good here or not
		domainTag := "domain:" + domain
		tags = append(tags, domainTag)

		// run the gauge
		err = c.Gauge("domain.expiration", daysLeft, tags, 1)
		if err != nil {
			log.Fatalf("Metric failed: %s %s %.0f", domain, tags, daysLeft)
		}
	}
	log.Infof("Finished! Domains checked: %d", len(domains))
}
