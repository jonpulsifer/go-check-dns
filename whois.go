package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/now"
	log "github.com/sirupsen/logrus"
)

// TODO: dry this up: make a common whois dial func / response scanner

func whois(server, domain string) *bufio.Scanner {
	if server == "" {
		server = "whois.iana.org"
	}
	log.Debugf("Finding %s on %s", domain, server)

	// https://tools.ietf.org/html/rfc3912
	conn, err := net.Dial("tcp", net.JoinHostPort(server, "43"))
	if err != nil {
		log.Fatalf("Could not connect to %s: %v", server, err)
	}
	//defer conn.Close()

	_, err = conn.Write((append([]byte(domain), '\r', '\n')))
	if err != nil {
		log.Fatalf("Write to server failed: %v", err)
	}
	log.Debugf("Queried %s for %s", server, domain)

	response := bufio.NewScanner(conn)
	response.Split(bufio.ScanLines)

	log.Debugf("Got WHOIS for %s, returning response", domain)
	return response
}

func findServer(domain string) string {
	var refer []string

	// get whois response
	response := whois("whois.iana.org", domain)

	log.Debugf("Scanning Response")

	for response.Scan() {
		// expected refer response:
		// refer: whois.example.com
		if strings.HasPrefix(response.Text(), "refer:") {
			refer = strings.Fields(response.Text())
			break
		}
	}
	if len(refer) == 0 {
		log.Fatalf("Could not get referer for %s from %s", domain, response.Text())
	}
	return refer[1]
}

func lookup(domain string) (time.Time, error) {
	// find the correct whois server
	refer := findServer(domain)
	log.WithFields(log.Fields{
		"refer":  refer,
		"domain": domain,
	}).Debugf("IANA refer: %s", refer)

	response := whois(refer, domain)

	r, err := regexp.Compile(".*[Ee]xpir.+")
	if err != nil {
		log.Fatalf("Could not compile regexp for expiration string")
	}

	// find expiry date
	var expiryDate []string
	for response.Scan() {
		if r.MatchString(response.Text()) {
			expiryDate = strings.Fields(response.Text())
			log.Debugf("Expiry date string: %s", response.Text())
			break
		}
	}

	if len(expiryDate) == 0 {
		err := errors.New("Could not find expiration date")
		return time.Now(), err
	}

	// more time formats
	now.TimeFormats = append(now.TimeFormats, "2006/01/02")           // .ca
	now.TimeFormats = append(now.TimeFormats, "02-Jan-2006")          // .com
	now.TimeFormats = append(now.TimeFormats, "2006-01-02T15:04:05Z") // .sucks

	// parse time
	log.Debugf("Attempting to parse: %s", expiryDate)
	log.Debugf("Time to parse: %s", expiryDate[len(expiryDate)-1])
	t, err := now.Parse(expiryDate[len(expiryDate)-1])
	if err != nil {
		err := fmt.Errorf("Could not parse time from %s: %v", expiryDate, err)
		return time.Now(), err
	}
	return t, nil
}
