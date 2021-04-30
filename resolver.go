package main

import (
	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/miekg/dns"
	"time"
)

func Resolve(name string, chaos bool, odohProxy string, upstream upstream.Upstream, rrTypes []uint16) ([]dns.RR, time.Duration, error) {
	var answers []dns.RR
	queryStartTime := time.Now()

	// Query for each requested RR type
	for _, qType := range rrTypes {
		// Create the DNS question
		req := dns.Msg{}

		if opts.DNSSEC {
			req.SetEdns0(4096, true)
		}

		// Set QCLASS
		var class uint16
		if chaos {
			class = dns.ClassCHAOS
		} else {
			class = dns.ClassINET
		}
		req.RecursionDesired = true
		req.Question = []dns.Question{{
			Name:   dns.Fqdn(name),
			Qtype:  qType,
			Qclass: class,
		}}

		var err error
		var reply *dns.Msg
		// Use upstream exchange if no ODoH proxy is configured
		if odohProxy == "" {
			// Send question to server
			reply, err = upstream.Exchange(&req)
		} else {
			reply, err = odohQuery(req, odohProxy, upstream.Address())
		}
		if err != nil {
			return nil, 0, err
		}

		answers = append(answers, reply.Answer...)
	}

	// Calculate total query time
	queryTime := time.Now().Sub(queryStartTime)

	return answers, queryTime, nil // nil error
}