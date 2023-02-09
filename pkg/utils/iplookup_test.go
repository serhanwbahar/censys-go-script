package utils_test

import (
	"censys-osint/pkg/utils"
	"testing"
)

// Tests to see if it can retrieve IPs of google.com
func TestIPv4Lookup(t *testing.T) {
	ips := utils.IPv4Lookup("google.com")
	if len(ips) == 0 {
		t.Error("No IPs were returned.")
	}
}
