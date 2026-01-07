// types/requests.go

package types

type ListRequest struct {
	MaxResults       int                `json:"maxResults,omitempty"`
	Filter           MarketFilter       `json:"filter,omitempty"`
	Sort             MarketSort         `json:"sort,omitempty"`
	MarketProjection []MarketProjection `json:"marketProjection,omitempty"`
}

type ListMarketBookRequest struct {
	MarketIds        []string           `json:"marketIds"`
	SelectionIds     []string           `json:"selectionIds,omitempty"`
	PriceProjection  PriceProjection    `json:"priceProjection,omitempty"`
	OrderProjection  OrderProjection    `json:"orderProjection,omitempty"`
	MarketProjection []MarketProjection `json:"marketProjection,omitempty"`
}
