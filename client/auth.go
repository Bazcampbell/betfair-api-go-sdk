// client/auth.go

package client

import (
	"betfair-api-go-sdk/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BASE_AUTH_URL = "https://identitysso-cert.betfair.com.au/api/"

// Returns session token + error
func (b *BetfairClient) login() (string, error) {
	loginUrl := BASE_AUTH_URL + "certlogin"

	params := url.Values{}
	params.Add("username", b.creds.Username)
	params.Add("password", b.creds.Password)
	paramsEncoded := params.Encode()

	req, err := http.NewRequest("POST", loginUrl, bytes.NewBufferString(paramsEncoded))
	if err != nil {
		return "", fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Set("X-Application", b.creds.AppKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %w", err)
	}

	var loginResp types.LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return "", fmt.Errorf("unable to parse response: %w", err)
	}

	if loginResp.Status == "SUCCESS" && len(loginResp.SessionToken) == 44 {
		return loginResp.SessionToken, nil
	}

	return "", fmt.Errorf("malformed response: %s, token: %s", loginResp.Status, loginResp.SessionToken)
}

// Returns error and assigns session token
func (b *BetfairClient) keepAlive() error {
	token, ok := b.sessionToken.Load().(string)
	if !ok || token == "" {
		return fmt.Errorf("session token not initialized")
	}

	keepAliveUrl := BASE_AUTH_URL + "keepAlive"

	req, err := http.NewRequest("POST", keepAliveUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Application", b.creds.AppKey)
	req.Header.Set("X-Authentication", token)

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %w", err)
	}

	var keepAliveResponse types.KeepAliveResponse
	err = json.Unmarshal(body, &keepAliveResponse)
	if err != nil {
		return fmt.Errorf("unable to parse response: %w", err)
	}

	if keepAliveResponse.Status == "SUCCESS" && len(keepAliveResponse.SessionToken) == 44 {
		b.sessionToken.Store(keepAliveResponse.SessionToken)
		return nil
	}

	return fmt.Errorf("malformed response: %s, error: %s, token: %s",
		keepAliveResponse.Status, keepAliveResponse.Error, keepAliveResponse.SessionToken,
	)
}

func (b *BetfairClient) logout() error {
	logoutUrl := BASE_AUTH_URL + "logout"

	req, err := http.NewRequest("POST", logoutUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	token, err := b.getSessionToken()
	if err != nil {
		return fmt.Errorf("error getting session token: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Application", b.creds.AppKey)
	req.Header.Set("X-Authentication", token)

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %w", err)
	}

	var logoutResponse types.LogoutResponse
	err = json.Unmarshal(body, &logoutResponse)
	if err != nil {
		return fmt.Errorf("unable to parse response: %w", err)
	}

	if logoutResponse.Status == "SUCCESS" {
		return nil
	}

	return fmt.Errorf("malformed response: %s, error: %s", logoutResponse.Status, logoutResponse.Error)
}
