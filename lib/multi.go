package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func GetMultiWhois(ctx context.Context, domains []string) ([]whoisparser.WhoisInfo, error) {
	allWhois := make([]whoisparser.WhoisInfo, 0, len(domains))

	for _, domain := range domains {
		var raw string
		var err error

		// Retry whois query up to 3 times
		for attempts := 0; attempts < 10; attempts++ {
			raw, err = whois.Whois(domain)

			if err == nil {
				break
			}

			// Wait before retrying
			time.Sleep(500 * time.Millisecond)
		}

		// Sleep to prevent rate limiting regardless of success
		time.Sleep(300 * time.Millisecond)

		if err != nil {
			fmt.Println("Domain:", domain, "Error:", err.Error())
			allWhois = append(allWhois, whoisparser.WhoisInfo{})
			continue
		}

		result, err := whoisparser.Parse(raw)
		if err != nil {
			fmt.Println("Domain:", domain, "Error: Not registered")
			allWhois = append(allWhois, whoisparser.WhoisInfo{})
			continue
		}

		// Check if Registrant is nil before accessing the Name field
		registrantName := "Unknown"
		if result.Registrant != nil {
			registrantName = result.Registrant.Name
		}

		fmt.Println("Domain:", domain, "Registrant:", registrantName)

		allWhois = append(allWhois, result)
	}

	return allWhois, nil
}
