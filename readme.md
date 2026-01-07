Betfair API Go SDK
==================

A lightweight, modern, production-ready Go client for the Betfair Exchange API (JSON-RPC).

Focuses on reliability, clean session management, automatic keep-alive, 
certificate authentication and graceful reconnection logic.

Features
--------
- Certificate-based login (base64 encoded cert + key)
- Automatic keep-alive with exponential backoff reconnect
- Generic JSON-RPC POST helper
- Implemented market discovery endpoints:
  - listEventTypes
  - listCompetitions
  - listCountries
  - listEvents
  - listMarketTypes
  - listMarketCatalogue
  - listMarketBook (with optional selectionIds filtering)
- Background error callback support
- Thread-safe
- Graceful shutdown

Installation
------------
```bash
    go get github.com/Bazcampbell/betfair-api-go-sdk
```

Quick Start Example
-------------------
```go
import (
	"betfair-api-go-sdk/client"
	"betfair-api-go-sdk/types"
	"fmt"
	"os"

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

    // All types implement stringer
	for _, c := range countries {
		fmt.Println(c.String())
	}
}

"""Prints:
    EventType: IdName{Id: 1, Name: Soccer}
    EventType: IdName{Id: 7522, Name: Basketball}
    EventType: IdName{Id: 2, Name: Tennis}
    EventType: IdName{Id: 3, Name: Golf}
    EventType: IdName{Id: 4, Name: Cricket}
    EventType: IdName{Id: 7524, Name: Ice Hockey}
    EventType: IdName{Id: 5, Name: Rugby Union}
    EventType: IdName{Id: 1477, Name: Rugby League}
    EventType: IdName{Id: 6, Name: Boxing}
    EventType: IdName{Id: 7, Name: Horse Racing}
    EventType: IdName{Id: 998917, Name: Volleyball}
    EventType: IdName{Id: 61420, Name: Australian Rules}
    EventType: IdName{Id: 136332, Name: Chess}
    EventType: IdName{Id: 3503, Name: Darts}
    EventType: IdName{Id: 26420387, Name: Mixed Martial Arts}
    EventType: IdName{Id: 4339, Name: Greyhound Racing}
    EventType: IdName{Id: 2378961, Name: Politics}
    EventType: IdName{Id: 6422, Name: Snooker}
    EventType: IdName{Id: 6423, Name: American Football}"""
```

Implemented Endpoints
---------------------
Market Discovery:
    ListEventTypes(filter)      → []ListEventTypesResponse
    ListCompetitions(filter)    → []ListCompetitionsResponse
    ListCountries(filter)       → []ListCountriesResponse
    ListEvents(filter)          → []ListEventsResponse
    ListMarketTypes(filter)     → []ListMarketTypesResponse

Market Data:
    ListMarketCatalogue(req)    → []ListMarketCataloguesResponse
    ListMarketBook(req)         → []ListMarketBookResponse
        (supports selectionIds filtering to reduce response size)

Fault Codes & Errors Reference
------------------------------
Official Betfair Cougar Fault Reporting Documentation:
https://betfair.github.io/cougar/legacy/Cougar_Fault_Reporting.html


Current Project Structure (as of Jan 2026)
-----------------------------------------
betfair-api-go-sdk/
├── client/
│   ├── client.go          # core client + keep-alive + lifecycle
│   ├── auth.go            # login/keepAlive/logout logic
│   └── list_endpoints.go  # all list*() market discovery methods
├── types/                 # Betfair request/response structs
└── util/                  # generic http/json helpers

To-Do / Missing (PRs welcome!)
------------------------------
- Trading endpoints: placeOrders, cancelOrders, updateOrders...
- Streaming API (Exchange Stream)
- Better structured error types (fault code mapping)
- Full context support
- Unit/integration tests

License
-------
MIT

Last updated: January 2026