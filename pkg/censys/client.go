package censys

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

// Censys Client with given base URL and authorization header.
type censysClient struct {
	BASE_URL      string
	AUTHORIZATION string
	SLEEP         time.Duration // Sleep a bit to avoid rate-limiting
	client        *http.Client
}

var CensysClient *censysClient = nil

// Initializes the Censys client. Should be called early in the program.
func Initalize(baseURL string, authToken string, sleepDuration time.Duration) error {
	if CensysClient != nil {
		return fmt.Errorf("censys: already initialized")
	}
	var cl censysClient

	cl.AUTHORIZATION = fmt.Sprintf("Basic %s", authToken)
	cl.BASE_URL = baseURL
	cl.SLEEP = sleepDuration
	cl.client = &http.Client{}

	CensysClient = &cl
	return nil
}

// Fetches a list of host names (in the form of a JSON array string) for the specified IP address.
func (censys *censysClient) GetHostNames(ip string) (string, error) {
	url := fmt.Sprintf("%s/api/v2/hosts/%s/names", censys.BASE_URL, ip)
	result, err := censys.genericRequest(url)
	if err != nil {
		return "", err
	}
	return result, nil
}

// Fetches the entire host entity by IP address and returns the
// most recent Censys view of the host and its services.
func (censys *censysClient) GetHost(ip string) (string, error) {
	url := fmt.Sprintf("%s/api/v2/hosts/%s", censys.BASE_URL, ip)
	result, err := censys.genericRequest(url)
	if err != nil {
		return "", err
	}
	return result, nil
}

// Fetches a list of events for the host with the specified IP address.
func (censys *censysClient) GetHostEvents(ip string) (string, error) {
	url := fmt.Sprintf("%s/api/v2/experimental/hosts/%s/events", censys.BASE_URL, ip)
	result, err := censys.genericRequest(url)
	if err != nil {
		return "", err
	}
	return result, nil
}

// Fetches the IPs of a host via Censys
func (censys *censysClient) GetHostsOfDomain(domain string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v2/hosts/search?q=%s", censys.BASE_URL, domain)
	result, err := censys.genericRequest(url)
	if err != nil {
		return []string{}, err
	}

	ipsJSON := gjson.Get(string(result), "hits.#.ip").String()

	var ips []string
	if err := json.Unmarshal([]byte(ipsJSON), &ips); err != nil {
		return []string{}, nil
	}
	return ips, nil
}

// A generic request for the Censys Search API, returns the JSON result.
// If the code is not 200 (OK), it will log the unexpected code.
func (censys *censysClient) genericRequest(url string) (string, error) {
	fmt.Println("\tRequesting at: ", url)
	time.Sleep(censys.SLEEP)

	// Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", censys.AUTHORIZATION)

	// Response
	res, err := censys.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	// Parse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	code := gjson.Get(string(body), "code").Int()
	if code != http.StatusOK {
		fmt.Println("\tUnexpected Code:", code)
	}
	return gjson.Get(string(body), "result").String(), nil
}

//// Note that a cleaner JSON approach like below is not feasible,
//// as JSON size is huge and varies greatly for each API endpoint
// type listHostNamesBody = struct {
// 	Code   int    `json:"code" bindings:"required"`
// 	Status string `json:"status"`
// 	Result struct {
// 		Ip    string   `json:"ip"`
// 		Names []string `json:"names"`
// 		Links struct {
// 			Next string `json:"next"`
// 		} `json:"links"`
// 	} `json:"result"`
// }
