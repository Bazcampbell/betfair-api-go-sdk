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

func TestClientInitFlow(t *testing.T) {
	_ = godotenv.Load("../.env")

	key := os.Getenv("BETFAIR_KEY_BASE64")
	cert := os.Getenv("BETFAIR_CERT_BASE64")
	username := os.Getenv("BETFAIR_USERNAME")
	password := os.Getenv("BETFAIR_PASSWORD")
	appKey := os.Getenv("BETFAIR_APP_KEY")
	proxy := os.Getenv("HTTP_PROXY")

	// Skip instead of fail if creds aren’t present
	if key == "" || cert == "" || username == "" || password == "" || appKey == "" {
		t.Fatalf("Betfair credentials not set")
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
		t.Fatalf("error creating session: %v", err)
	}

	filter := types.MarketFilter{
		MarketCountries: []string{"AU", "GB", "US"},
	}

	eventTypes, err := bfClient.ListEventTypes(filter)
	if err != nil {
		t.Fatalf("error listing event types: %v", err)
	}

	if len(eventTypes) == 0 {
		t.Fatal("expected event types, got none")
	}

	// Build lookup map
	found := make(map[string]bool)
	for _, e := range eventTypes {
		found[e.EventType.Name] = true
	}

	required := []string{
		"Soccer",
		"Tennis",
		"Horse Racing",
		"Cricket",
		"Australian Rules",
	}

	for _, name := range required {
		if !found[name] {
			t.Fatalf("expected event type %q to exist", name)
		}
	}
}

func TestClientUsesGoodProxy(t *testing.T) {
	_ = godotenv.Load("../.env")

	proxy := os.Getenv("HTTP_PROXY")
	if proxy == "" {
		t.Fatalf("HTTP_PROXY not set")
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

func TestClientUsesBadProxy(t *testing.T) {
	_ = godotenv.Load("../.env")

	proxy := os.Getenv("HTTP_PROXY") + "badport"
	if proxy == "" {
		t.Fatalf("HTTP_PROXY not set")
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
		ProxyUrl:   &proxy,
	}

	onErrorFunc := func(err error) {
		t.Fatalf("betfair client error: %v", err)
	}

	_, err := client.NewSession(creds, onErrorFunc)
	if err == nil {
		t.Fatalf("uncaught bad proxy")
	}

	if !strings.Contains(err.Error(), "invalid port") {
		t.Fatalf("uncaught bad proxy")
	}
}

func TestClientUsesNonAusProxy(t *testing.T) {
	_ = godotenv.Load("../.env")

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
