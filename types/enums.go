// types/enums.go

package types

type MarketSort string

const (
	MINIMUM_TRADED    MarketSort = "MINIMUM_TRADED"
	MAXIMUM_TRADED    MarketSort = "MAXIMUM_TRADED"
	MINIMUM_AVAILABLE MarketSort = "MINIMUM_AVAILABLE"
	MAXIMUM_AVAILABLE MarketSort = "MAXIMUM_AVAILABLE"
	FIRST_TO_START    MarketSort = "FIRST_TO_START"
	LAST_TO_START     MarketSort = "LAST_TO_START"
)

type MarketProjection string

const (
	COMPETITION        MarketProjection = "COMPETITION"
	EVENT              MarketProjection = "EVENT"
	EVENT_TYPE         MarketProjection = "EVENT_TYPE"
	MARKET_START_TIME  MarketProjection = "MARKET_START_TIME"
	MARKET_DESCRIPTION MarketProjection = "MARKET_DESCRIPTION"
	RUNNER_DESCRIPTION MarketProjection = "RUNNER_DESCRIPTION"
	RUNNER_METADATA    MarketProjection = "RUNNER_METADATA"
)

type PriceData string

const (
	SP_AVAILABLE   PriceData = "SP_AVAILABLE"
	SP_TRADED      PriceData = "SP_TRADED"
	EX_BEST_OFFERS PriceData = "EX_BEST_OFFERS"
	EX_ALL_OFFERS  PriceData = "EX_ALL_OFFERS"
	EX_TRADED      PriceData = "EX_TRADED"
)

type OrderProjection string

const (
	ALL                OrderProjection = "ALL"
	EXECUTABLE         OrderProjection = "EXECUTABLE"
	EXECUTION_COMPLETE OrderProjection = "EXECUTION_COMPLETE"
)

type MarketStatus string

const (
	INACTIVE  MarketStatus = "INACTIVE"
	OPEN      MarketStatus = "OPEN"
	SUSPENDED MarketStatus = "SUSPENDED"
	CLOSED    MarketStatus = "CLOSED"
)

type Side string

const (
	BACK Side = "BACK"
	LAY  Side = "LAY"
)

type MatchProjection string

const (
	NO_ROLLUP              MatchProjection = "NO_ROLLUP"
	ROLLED_UP_BY_PRICE     MatchProjection = "ROLLED_UP_BY_PRICE"
	ROLLED_UP_BY_AVG_PRICE MatchProjection = "ROLLED_UP_BY_AVG_PRICE"
)

type RollupModel string

const (
	STAKE  RollupModel = "STAKE"
	PAYOUT RollupModel = "PAYOUT"
	NONE   RollupModel = "NONE"
)
