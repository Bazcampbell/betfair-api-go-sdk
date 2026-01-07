// util/http.go

package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	maxRetries = 3
	baseDelay  = time.Millisecond * 800
	BASE_URL   = "https://api.betfair.com/exchange/betting/rest/v1.0/"
)

// Send a POST request to a given URL
// Add parameter headers
// Attempt to unmarshal the response into T
func GenericPost[T any](client *http.Client, endpoint, appKey, sessionToken string, body any) (T, error) {
	var result T
	var lastErr error

	fullUrl := BASE_URL + endpoint
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// exponential backoff + jitter
			delay := baseDelay * time.Duration(1<<attempt)
			jitter := time.Duration(time.Now().UnixNano()%100) * time.Millisecond
			time.Sleep(delay + jitter)
		}

		reqBody, err := json.Marshal(body)
		if err != nil {
			return result, fmt.Errorf("unable to marshal body: %w", err)
		}

		req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(reqBody))
		if err != nil {
			return result, fmt.Errorf("unable to build request: %w", err)
		}

		req.Header.Set("X-Application", appKey)
		req.Header.Set("X-Authentication", sessionToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("unable to make request: %w", err)
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("unable to read response: %w", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return result, fmt.Errorf("http %d: %s", resp.StatusCode, string(bodyBytes))
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("status code: %d, %s", resp.StatusCode, body)

			if !shouldRetry(resp.StatusCode, attempt) {
				break
			}
			continue
		}

		// retry certain statuses
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("status code: %d, %s", resp.StatusCode, body)

			if !shouldRetry(resp.StatusCode, attempt) {
				break
			}
			continue
		}

		if err = json.Unmarshal(bodyBytes, &result); err != nil {
			lastErr = fmt.Errorf("json unmarshal failed: %w", err)
			continue
		}

		return result, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("request failed after %d attempts (unknown reason)", maxRetries)
	}

	return result, fmt.Errorf("%w (after %d attempts)", lastErr, maxRetries)
}

func shouldRetry(status int, attempt int) bool {
	if attempt >= maxRetries-1 {
		return false // last attempt anyway
	}

	// always retry these
	if status == 429 || (status >= 500 && status <= 599) {
		return true
	}

	// usually don't retry these other status codes
	return false
}

// Helper function to parse comma-separated query params into string slice
// Returns nil if parameter is not present or empty
func ParseQueryArrayOrNil(r *http.Request, param string) []string {
	value := r.URL.Query().Get(param)
	if value == "" {
		return nil
	}

	// Split by comma
	var result []string
	for _, v := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	// Return nil if no valid values found
	if len(result) == 0 {
		return nil
	}

	return result
}
