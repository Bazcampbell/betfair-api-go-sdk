// client/client.go

package client

// Handles state and authentication.

import (
	"betfair-api-go-sdk/types"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type BetfairClient struct {
	client       *http.Client
	sessionToken atomic.Value
	creds        types.BetfairCredentials

	mu     sync.RWMutex
	wg     sync.WaitGroup
	closed atomic.Bool
	ctx    context.Context
	cancel context.CancelFunc

	onError func(error)
}

const (
	keepAliveRetries    = 3
	keepAliveRetryDelay = 1 * time.Second
)

func NewSession(creds types.BetfairCredentials, onError func(error)) (*BetfairClient, error) {
	if creds.AppKey == "" {
		return nil, fmt.Errorf("app key cannot be empty")
	}

	if creds.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if creds.Password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	certPEM, err := base64.StdEncoding.DecodeString(creds.CertString)
	if err != nil {
		return nil, fmt.Errorf("invalid cert string. certificate string must be base64 encoded: %w", err)
	}

	keyPEM, err := base64.StdEncoding.DecodeString(creds.KeyString)
	if err != nil {
		return nil, fmt.Errorf("invalid key string. key string must be base64 encoded: %w", err)
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert/key pair: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if creds.ProxyUrl != nil {
		proxyUrl, err := url.Parse(*creds.ProxyUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}

		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	client := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())

	b := &BetfairClient{
		client:  &client,
		onError: onError,
		creds:   creds,
		ctx:     ctx,
		cancel:  cancel,
	}

	b.keepAliveTicker()

	sessionToken, err := b.login()
	if err != nil {
		return nil, fmt.Errorf("unable to login: %w", err)
	}

	b.sessionToken.Store(sessionToken)

	return b, nil
}

func (b *BetfairClient) keepAliveTicker() {
	b.wg.Add(1)
	ticker := time.NewTicker(6 * time.Hour)
	consecutiveFailures := 0

	go func() {
		defer b.wg.Done()
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if b.closed.Load() {
					return
				}

				if err := b.keepAlive(); err != nil {
					consecutiveFailures++
					b.onError(fmt.Errorf("keepAlive failed: %w", err))

					// Exponential backoff before reconnect attempt
					if consecutiveFailures > 1 {
						delay := keepAliveRetryDelay * time.Duration(1<<(consecutiveFailures-1))
						jitter := time.Duration(rand.Int63n(100)) * time.Millisecond
						time.Sleep(delay + jitter)
					}

					// Try to reconnect
					if reconnectErr := b.reconnect(); reconnectErr != nil {
						if consecutiveFailures >= keepAliveRetries {
							b.onError(fmt.Errorf("max reconnect attempts reached, closing client"))
							if err := b.logout(); err != nil {
								b.onError(fmt.Errorf("unable to logout: %w", err))
							}
							b.close()
							return
						}
					} else {
						consecutiveFailures = 0
					}
				} else {
					consecutiveFailures = 0
				}

			case <-b.ctx.Done():
				if err := b.logout(); err != nil {
					b.onError(fmt.Errorf("logout failed: %w", err))
				}
				return
			}
		}
	}()
}

func (b *BetfairClient) getSessionToken() (string, error) {
	val := b.sessionToken.Load()
	if val == nil {
		return "", fmt.Errorf("token not initialised")
	}

	token, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("invalid token type")
	}

	return token, nil
}

func (b *BetfairClient) reconnect() error {
	token, err := b.login()
	if err != nil {
		return err
	}

	b.sessionToken.Store(token)
	return nil
}

func (b *BetfairClient) close() error {
	if !b.closed.CompareAndSwap(false, true) {
		return fmt.Errorf("client already closed")
	}

	b.cancel()

	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Clean shutdown
		return nil
	case <-time.After(10 * time.Second):
		// Timeout waiting for goroutines
		return fmt.Errorf("timeout waiting for client shutdown")
	}
}
