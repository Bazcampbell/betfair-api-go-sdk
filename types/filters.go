// types/filters.go

package types

import "time"

type MarketFilter struct {
	EventTypeIds    []string   `json:"eventTypeIds,omitempty"`
	EventIds        []string   `json:"eventIds,omitempty"`
	ComptetitionIds []string   `json:"competitionIds,omitempty"`
	MarketIds       []string   `json:"marketIds,omitempty"`
	MarketCountries []string   `json:"marketCountries,omitempty"`
	MarketTypeCodes []string   `json:"marketTypeCodes,omitempty"`
	TextQuery       string     `json:"textQuery,omitempty"`
	TimeRange       *TimeRange `json:"marketStartTime,omitempty"`
}

type PriceProjection struct {
	PriceData  []PriceData           `json:"priceData,omitempty"`
	Overrides  ExBestOffersOverrides `json:"exBestOffersOverrides,omitempty"`
	Virtualise bool                  `json:"virtualise,omitempty"`
}

type ExBestOffersOverrides struct {
	BestPricesDepth int         `json:"bestPricesDepth,omitempty"`
	RollupModel     RollupModel `json:"rollupModel,omitempty"`
	RollupLimit     int         `json:"rollupLimit,omitempty"` // all prices with stake above this are returned
}

type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
