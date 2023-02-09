package censys_test

import (
	"censys-osint/internal/constants"
	"censys-osint/pkg/censys"
	"censys-osint/pkg/utils"
	"encoding/base64"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var ip string
var authToken string

func TestMain(m *testing.M) {
	// Parse env variables to obtain auth token
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Could not read .env, see README.md")
	}
	apiId, apiSecret := os.Getenv(constants.ENV_CENSYS_API_ID), os.Getenv(constants.ENV_CENSYS_API_SECRET)
	authToken = base64.URLEncoding.EncodeToString([]byte(apiId + ":" + apiSecret)) // ':' is important

	// Use google ip as default
	ip = utils.IPv4Lookup("google.com")[0]

	m.Run()
}

func TestCensys(t *testing.T) {
	// Initialize
	if err := censys.Initalize(constants.CENSYS_BASE_URL, authToken, constants.CENSYS_REQ_SLEEP); err != nil {
		t.Error(err)
	}
	if err := censys.Initalize(constants.CENSYS_BASE_URL, authToken, constants.CENSYS_REQ_SLEEP); err == nil {
		t.Error("Expected to fail on second initialize")
	}

	censysClient := censys.CensysClient

	if _, err := censysClient.GetHost(ip); err != nil {
		t.Error(err)
	}
	if _, err := censysClient.GetHostNames(ip); err != nil {
		t.Error(err)
	}
	if _, err := censysClient.GetHostEvents(ip); err != nil {
		t.Error(err)
	}

}
