// types/formatters.go

package types

import "fmt"

// Event String method
func (e Event) String() string {
	return fmt.Sprintf("Event{Id: %s, Name: %s, Country: %s, Timezone: %s, OpenDate: %s}",
		e.Id, e.Name, e.CountryCode, e.Timezone, e.OpenDate)
}

// IdName String method
func (i IdName) String() string {
	return fmt.Sprintf("IdName{Id: %s, Name: %s}", i.Id, i.Name)
}

// ListMarketTypesResponse String method
func (l ListMarketTypesResponse) String() string {
	return fmt.Sprintf("MarketType: %s", l.MarketType)
}

// ListEventTypesResponse String method
func (l ListEventTypesResponse) String() string {
	return fmt.Sprintf("EventType: %s", l.EventType.String())
}

// ListCompetitionsResponse String method
func (l ListCompetitionsResponse) String() string {
	return fmt.Sprintf("Competition: %s, Region: %s", l.Competition.String(), l.Region)
}

// ListCountriesResponse String method
func (l ListCountriesResponse) String() string {
	return fmt.Sprintf("CountryCode: %s", l.CountryCode)
}

// ListEventsResponse String method
func (l ListEventsResponse) String() string {
	return fmt.Sprintf("Event: %s", l.Event.String())
}

// ListMarketCataloguesResponse String method
func (l ListMarketCataloguesResponse) String() string {
	result := fmt.Sprintf("Market{Id: %s, Name: %s, StartTime: %s, TotalMatched: %.2f}\n",
		l.MarketId, l.MarketName, l.MarketStartTime, l.TotalMatched)

	for i, runner := range l.Runners {
		result += fmt.Sprintf("  Runner %d: %s\n", i+1, runner.String())
	}

	return result
}

// Runner String method
func (r Runner) String() string {
	name := r.RunnerName
	if name == "" {
		name = "<no name>"
	}
	return fmt.Sprintf("Runner{Id: %d, Name: %s, Handicap: %.1f, TotalMatched: %.2f, Back: %d, Lay: %d, Traded: %d}",
		r.SelectionId, name, r.Handicap, r.TotalMatched, len(r.Ex.Back), len(r.Ex.Lay), len(r.Ex.Traded))
}

// Ex String method
func (e Ex) String() string {
	result := "Ex{\n"

	if len(e.Back) > 0 {
		result += "    Back: "
		for _, price := range e.Back {
			result += fmt.Sprintf("[%s] ", price.String())
		}
		result += "\n"
	}

	if len(e.Lay) > 0 {
		result += "    Lay: "
		for _, price := range e.Lay {
			result += fmt.Sprintf("[%s] ", price.String())
		}
		result += "\n"
	}

	if len(e.Traded) > 0 {
		result += "    Traded: "
		for _, price := range e.Traded {
			result += fmt.Sprintf("[%s] ", price.String())
		}
		result += "\n"
	}

	result += "  }"
	return result
}

// RunnerPrice String method
func (r RunnerPrice) String() string {
	return fmt.Sprintf("%.2f@%.2f", r.Price, r.Size)
}

// ListMarketBookResponse String method
func (l ListMarketBookResponse) String() string {
	result := fmt.Sprintf("MarketBook{%d Runners}\n", len(l.Runners))

	for i, runner := range l.Runners {
		result += fmt.Sprintf("\n  Runner %d - SelectionId: %d, Handicap: %.1f, TotalMatched: %.2f\n",
			i+1, runner.SelectionId, runner.Handicap, runner.TotalMatched)

		if len(runner.Ex.Back) > 0 {
			result += "    Back: "
			for _, price := range runner.Ex.Back {
				result += fmt.Sprintf("[%.2f @ %.2f] ", price.Price, price.Size)
			}
			result += "\n"
		} else {
			result += "    Back: [none]\n"
		}

		if len(runner.Ex.Lay) > 0 {
			result += "    Lay:  "
			for _, price := range runner.Ex.Lay {
				result += fmt.Sprintf("[%.2f @ %.2f] ", price.Price, price.Size)
			}
			result += "\n"
		} else {
			result += "    Lay:  [none]\n"
		}

		if len(runner.Ex.Traded) > 0 {
			result += "    Traded: "
			for _, price := range runner.Ex.Traded {
				result += fmt.Sprintf("[%.2f @ %.2f] ", price.Price, price.Size)
			}
			result += "\n"
		}
	}

	return result
}

// LoginResponse String method
func (l LoginResponse) String() string {
	return fmt.Sprintf("Login{Status: %s, SessionToken: %s...}", l.Status, truncate(l.SessionToken, 10))
}

// KeepAliveResponse String method
func (k KeepAliveResponse) String() string {
	if k.Error != "" {
		return fmt.Sprintf("KeepAlive{Status: %s, Error: %s}", k.Status, k.Error)
	}
	return fmt.Sprintf("KeepAlive{Status: %s, Token: %s...}", k.Status, truncate(k.SessionToken, 10))
}

// LogoutResponse String method
func (l LogoutResponse) String() string {
	if l.Error != "" {
		return fmt.Sprintf("Logout{Status: %s, Error: %s}", l.Status, l.Error)
	}
	return fmt.Sprintf("Logout{Status: %s}", l.Status)
}

// Helper function to truncate strings
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
