package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	dns "google.golang.org/api/dns/v1"
)

func getZones(project string) []string {
	// return this
	var zones []string

	// google cloud api client
	// WARNING: this uses the cloud-platform-ro scope
	// use IAM to restrict api access
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, dns.CloudPlatformReadOnlyScope)
	if err != nil {
		log.Fatalf("Could not create Google API client: %v", err)
	}

	// new DNS api client
	dnsService, err := dns.New(client)
	if err != nil {
		log.Fatalf("Could not create DNS client: %v", err)
	}

	// gets managed zones from a GCP project
	resp, err := dnsService.ManagedZones.List(project).Do()
	if err != nil {
		log.Fatalf("Could not list ManagedZones in %s: %v", project, err)
	}

	// iterate over the list
	for _, zone := range resp.ManagedZones {
		if len(strings.Split(zone.DnsName, ".")) == 3 {
			zones = append(zones, strings.TrimRight(zone.DnsName, "."))
			log.Debugf("Found %s", zone.DnsName)
		} else {
			log.Debugf("Skipping subdomain: %s", zone.DnsName)
		}
	}
	return zones
}
