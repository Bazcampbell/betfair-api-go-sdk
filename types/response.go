// types/response.go

package types

type Event struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
	Timezone    string `json:"timezone"`
	OpenDate    string `json:"openDate"`
}

type IdName struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListMarketTypesResponse struct {
	MarketType  string `json:"marketType"`
	MarketCount int    `json:"marketCount"`
}

type ListEventTypesResponse struct {
	EventType   IdName `json:"eventType"`
	MarketCount int    `json:"marketCount"`
}

type ListCompetitionsResponse struct {
	Competition IdName `json:"competition"`
	Region      string `json:"competitionRegion"`
	MarketCount int    `json:"marketCount"`
}

type ListCountriesResponse struct {
	CountryCode string `json:"countryCode"`
	MarketCount int    `json:"marketCount"`
}

type ListEventsResponse struct {
	Event       Event `json:"event"`
	MarketCount int   `json:"marketCount"`
}

type ListMarketCataloguesResponse struct {
	MarketId        string   `json:"marketId"`
	MarketName      string   `json:"marketName"`
	MarketStartTime string   `json:"marketStartTime,omitempty"`
	TotalMatched    float32  `json:"totalMatched"`
	Runners         []Runner `json:"runners,omitempty"`
}

type ListMarketSelectionsResponse struct {
}

type Runner struct {
	SelectionId     int     `json:"selectionId"`
	RunnerName      string  `json:"runnerName"`
	Handicap        float32 `json:"handicap"`
	LastPriceTraded float32 `json:"lastPriceTraded"`
	TotalMatched    float32 `json:"totalMatched"`
	Ex              Ex      `json:"ex"`
}

type Ex struct {
	Back   []RunnerPrice `json:"availableToBack"`
	Lay    []RunnerPrice `json:"availableToLay"`
	Traded []RunnerPrice `json:"tradedVolume"`
}

type RunnerPrice struct {
	Price float32 `json:"price"`
	Size  float32 `json:"size"`
}

type ListMarketBookResponse struct {
	Runners []Runner `json:"runners"`
}

// AUTH RESPONSE TYPES
type LoginResponse struct {
	SessionToken string `json:"sessionToken"`
	Status       string `json:"loginStatus"`
}

type KeepAliveResponse struct {
	SessionToken string `json:"token"`
	Status       string `json:"status"`
	Error        string `json:"error"`
}

type LogoutResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
