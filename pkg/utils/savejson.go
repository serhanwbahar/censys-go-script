package utils

import (
	"fmt"
	"log"
	"os"
)

// Saves the given ipResults for a domain at the specified filename.
// Format is in JSON.
func SaveJSON(filename string, domain string, ipResults string) {
	jsonResult := fmt.Sprintf(`{"domain": "%s", "results": %s}`, domain, ipResults)
	fmt.Println("\nWriting results to file: ", filename)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(jsonResult); err != nil {
		log.Fatal(err)
	}
}
