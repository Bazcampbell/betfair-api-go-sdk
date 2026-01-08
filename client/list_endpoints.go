// client/list_endpoints.go

package client

import (
	"betfair-api-go-sdk/types"
	"betfair-api-go-sdk/util"
	"fmt"
)

const BASE_URL = "https://api.betfair.com/exchange/betting/rest/v1.0/"

func (b *BetfairClient) ListEventTypes(filter types.MarketFilter) ([]types.ListEventTypesResponse, error) {
	token, err := b.getSessionToken()
	if err != nil {
		return nil, err
	}

	body := types.ListRequest{Filter: filter}
	return util.GenericPost[[]types.ListEventTypesResponse](b.client, "listEventTypes/", b.creds.AppKey, token, body)
}

func (b *BetfairClient) ListCompetitions(filter types.MarketFilter) ([]types.ListCompetitionsResponse, error) {
	token, err := b.getSessionToken()
	if err != nil {
		return nil, err
	}

	body := types.ListRequest{Filter: filter}
	return util.GenericPost[[]types.ListCompetitionsResponse](b.client, "listCompetitions/", b.creds.AppKey, token, body)
}

func (b *BetfairClient) ListCountries(filter types.MarketFilter) ([]types.ListCountriesResponse, error) {
	token, err := b.getSessionToken()
	if err != nil {
		return nil, err
	}

	body := types.ListRequest{Filter: filter}
	return util.GenericPost[[]types.ListCountriesResponse](b.client, "listCountries/", b.creds.AppKey, token, body)
}

func (b *BetfairClient) ListEvents(filter types.MarketFilter) ([]types.ListEventsResponse, error) {
	token, err := b.getSessionToken()
	if err != nil {
		return nil, err
	}

	body := types.ListRequest{Filter: filter}
	return util.GenericPost[[]types.ListEventsResponse](b.client, "listEvents/", b.creds.AppKey, token, body)
}

func (b *BetfairClient) ListMarketTypes(filter types.MarketFilter) ([]types.ListMarketTypesResponse, error) {
	token, err := b.getSessionToken()
	if err != nil {
		return nil, err
	}

	body := types.ListRequest{Filter: filter}
	return util.GenericPost[[]types.ListMarketTypesResponse](b.client, "listMarketTypes/", b.creds.AppKey, token, body)
}

func (b *BetfairClient) ListMarketCatalogues(req types.ListRequest) ([]types.ListMarketCataloguesResponse, error) {
	token, err := b.getSessionToken()
	if err != nil {
		return nil, err
	}

	return util.GenericPost[[]types.ListMarketCataloguesResponse](b.client, "listMarketCatalogue/", b.creds.AppKey, token, req)
}

func (b *BetfairClient) ListMarketBook(req types.ListMarketBookRequest) ([]types.ListMarketBookResponse, error) {
	token, err := b.getSessionToken()
	if err != nil {
		return nil, err
	}

	result, err := util.GenericPost[[]types.ListMarketBookResponse](b.client, "listMarketBook/", b.creds.AppKey, token, req)
	if err != nil {
		return nil, err
	}

	// Filter runners to only include requested selectionIds
	if len(req.SelectionIds) > 0 {
		selectionSet := make(map[string]bool)
		for _, id := range req.SelectionIds {
			selectionSet[id] = true
		}

		for i := range result {
			var filteredRunners []types.Runner
			for _, runner := range result[i].Runners {
				if selectionSet[fmt.Sprintf("%d", runner.SelectionId)] {
					filteredRunners = append(filteredRunners, runner)
				}
			}
			result[i].Runners = filteredRunners
		}
	}

	return result, nil
}
