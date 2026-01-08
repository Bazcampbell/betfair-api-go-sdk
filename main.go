// main.go

package main

import (
	"fmt"
	"os"

	"github.com/Bazcampbell/betfair-api-go-sdk/client"
	"github.com/Bazcampbell/betfair-api-go-sdk/types"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	key := os.Getenv("BETFAIR_KEY_BASE64")
	cert := os.Getenv("BETFAIR_CERT_BASE64")
	username := os.Getenv("BETFAIR_USERNAME")
	password := os.Getenv("BETFAIR_PASSWORD")
	appKey := os.Getenv("BETFAIR_APP_KEY")

	creds := types.BetfairCredentials{
		Username:   username,
		Password:   password,
		AppKey:     appKey,
		KeyString:  key,
		CertString: cert,
	}

	onErrorFunc := func(err error) {
		fmt.Println("betfair client error: %w", err)
	}

	client, err := client.NewSession(creds, onErrorFunc)
	if err != nil {
		fmt.Println("error creating session: %w", err)
		return
	}

	filter := types.MarketFilter{
		MarketCountries: []string{"AU", "GB", "US"},
	}

	countries, err := client.ListEventTypes(filter)
	if err != nil {
		fmt.Println("error listing event types: %w", err)
		return
	}

	for _, c := range countries {
		fmt.Println(c.String())
	}
}
