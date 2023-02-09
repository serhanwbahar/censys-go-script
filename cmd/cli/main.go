package main

import (
	"censys-osint/internal/constants"
	"censys-osint/pkg/censys"
	"censys-osint/pkg/utils"
	"encoding/base64"

	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Parse command line args
	args := os.Args
	argc := len(os.Args)
	if argc != 2 {
		log.Fatal("Expected a command line arg")
	}
	domain := args[1]
	fmt.Println("Domain: ", domain)
	// It will be saved to this file
	filename := fmt.Sprintf("./out/%s.json", domain)

	// Parse env variables to obtain auth token
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Could not read .env, see README.md")
	}
	apiId, apiSecret := os.Getenv(constants.ENV_CENSYS_API_ID), os.Getenv(constants.ENV_CENSYS_API_SECRET)
	if apiId == "" {
		log.Fatal("API ID is missing, see README.md")
	}
	if apiSecret == "" {
		log.Fatal("API SECRET is missing, see README.md")
	}
	authToken := base64.URLEncoding.EncodeToString([]byte(apiId + ":" + apiSecret)) // ':' is important

	// Initialize Censys client
	if err := censys.Initalize(constants.CENSYS_BASE_URL, authToken, constants.CENSYS_REQ_SLEEP); err != nil {
		log.Fatal("Could not initialize Censys client")
	}
	censysClient := censys.CensysClient

	var ips []string
	if constants.USE_LOCAL_IPLOOKUP {
		// Use local ipv4 lookup to fetch IPs, will return a lot less amount than Censys
		ips = utils.IPv4Lookup(domain)
	} else {
		// Use Censys to fetch IPs, will return a lot of IPs. We only take first at most 50 ips for now.
		// i.e. no pagination continued
		if hostIps, err := censysClient.GetHostsOfDomain(domain); err != nil {
			log.Fatal(err)
		} else {
			ips = hostIps
		}
	}

	// Fetch data for each found ip
	fmt.Printf("Found %d IPs:\n", len(ips))
	ipResults := "["
	for i, ip := range ips {
		fmt.Printf("\n\t[%d] IPv4: %s\n", i, ip)

		host, err := censysClient.GetHost(ip)
		if err != nil {
			log.Fatal(err)
		}
		hostNames, err := censysClient.GetHostNames(ip)
		if err != nil {
			log.Fatal(err)
		}
		hostEvents, err := censysClient.GetHostEvents(ip)
		if err != nil {
			log.Fatal(err)
		}

		// Store results as JSON strings
		ipResults += fmt.Sprintf(`{"ip":"%s","host": %s,"host_names": %s,"host_events": %s}`, ip, host, hostNames, hostEvents)

		// Save intermediary results for performance
		if i > 0 && i%constants.SAVE_TO_FILE_INTERVAL == 0 {
			fmt.Printf("\n[%d/%d]Saving intermediary result at %s\n", i, len(ips), filename)
			utils.SaveJSON(filename, domain, ipResults+"]")
		}

		if i != len(ips)-1 {
			ipResults += ","
		}
	}

	fmt.Println("\nWriting final results at ", filename)
	utils.SaveJSON(filename, domain, ipResults+"]")

}
