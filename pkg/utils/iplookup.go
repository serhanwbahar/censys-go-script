package utils

import (
	"net"
)

// Returns the list of IPv4s of a given domain.
func IPv4Lookup(domain string) []string {
	ips, _ := net.LookupIP(domain) // e.g. google.com
	ans := make([]string, 0)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			ans = append(ans, ipv4.String())
		}
	}
	return ans
}
