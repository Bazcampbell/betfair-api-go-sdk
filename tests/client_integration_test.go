//go:build integration

package tests

import (
	"os"
	"strings"
	"testing"

	"betfair-api-go-sdk/client"
	"betfair-api-go-sdk/types"

	"github.com/joho/godotenv"
)

func TestClientUsesGoodProxy(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("unable to load .env file: %s", err)
	}

	proxy := os.Getenv("HTTP_PROXY")
	key := os.Getenv("BETFAIR_KEY_BASE64")
	cert := os.Getenv("BETFAIR_CERT_BASE64")
	username := os.Getenv("BETFAIR_USERNAME")
	password := os.Getenv("BETFAIR_PASSWORD")
	appKey := os.Getenv("BETFAIR_APP_KEY")

	if key == "" || cert == "" || username == "" || password == "" || appKey == "" || proxy == "" {
		t.Fatalf("Betfair credentials or proxy not set")
	}

	creds := types.BetfairCredentials{
		Username:   username,
		Password:   password,
		AppKey:     appKey,
		KeyString:  key,
		CertString: cert,
		ProxyUrl:   &proxy,
	}

	onErrorFunc := func(err error) {
		t.Fatalf("betfair client error: %v", err)
	}

	bfClient, err := client.NewSession(creds, onErrorFunc)
	if err != nil {
		t.Fatalf("error creating session with proxy: %v", err)
	}

	_, err = bfClient.ListEventTypes(types.MarketFilter{
		MarketCountries: []string{"AU"},
	})
	if err != nil {
		t.Fatalf("betfair call failed via proxy: %v", err)
	}
}

func TestClientUsesNonAusProxy(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("unable to load .env file: %s", err)
	}

	key := os.Getenv("BETFAIR_KEY_BASE64")
	cert := os.Getenv("BETFAIR_CERT_BASE64")
	username := os.Getenv("BETFAIR_USERNAME")
	password := os.Getenv("BETFAIR_PASSWORD")
	appKey := os.Getenv("BETFAIR_APP_KEY")

	if key == "" || cert == "" || username == "" || password == "" || appKey == "" {
		t.Fatalf("Betfair credentials not set")
	}

	creds := types.BetfairCredentials{
		Username:   username,
		Password:   password,
		AppKey:     appKey,
		KeyString:  key,
		CertString: cert,
	}

	onErrorFunc := func(err error) {
		t.Fatalf("betfair client error: %v", err)
	}

	_, err := client.NewSession(creds, onErrorFunc)
	if err == nil {
		t.Fatal("expected BETTING_RESTRICTED_LOCATION error, got none")
	}

	if !strings.Contains(err.Error(), "BETTING_RESTRICTED_LOCATION") {
		t.Fatalf("expected BETTING_RESTRICTED_LOCATION error, got: %v", err)
	}
}
