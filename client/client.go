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
	"sync"
	"sync/atomic"
	"time"
)

type BetfairClient struct {
	client       *http.Client
	sessionToken string
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
	certPEM, err := base64.StdEncoding.DecodeString(creds.CertString)
	if err != nil {
		return nil, fmt.Errorf("invalid cert string. certificate string must be base64 encoded: %w", err)
	}

	keyPEM, err := base64.StdEncoding.DecodeString(creds.KeyString)
	if err != nil {
		return nil, fmt.Errorf("invalid cert string. key string must be base64 encoded: %w", err)
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

	b.sessionToken = sessionToken

	return b, nil
}

// OnError registers a callback for background errors
// Pass nil to clear the handler
func (b *BetfairClient) OnError(handler func(error)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.onError = handler
}

func (b *BetfairClient) notifyError(err error) {
	b.mu.RLock()
	handler := b.onError
	b.mu.RUnlock()

	if handler != nil {
		// Non-blocking call
		go handler(err)
	}
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
					b.notifyError(fmt.Errorf("keepAlive failed: %w", err))

					// Exponential backoff before reconnect attempt
					if consecutiveFailures > 1 {
						delay := keepAliveRetryDelay * time.Duration(1<<(consecutiveFailures-1))
						jitter := time.Duration(rand.Int63n(100)) * time.Millisecond
						time.Sleep(delay + jitter)
					}

					// Try to reconnect
					if reconnectErr := b.reconnect(); reconnectErr != nil {
						b.notifyError(fmt.Errorf("reconnect failed (attempt %d/%d): %w",
							consecutiveFailures, keepAliveRetries, reconnectErr))

						if consecutiveFailures >= keepAliveRetries {
							b.notifyError(fmt.Errorf("max reconnect attempts reached, closing client"))
							if err := b.logout(); err != nil {
								b.notifyError(fmt.Errorf("unable to logout: %w", err))
							}
							b.close()
							return
						}
					} else {
						// Success! Reset counter
						consecutiveFailures = 0
					}
				} else {
					consecutiveFailures = 0
				}

			case <-b.ctx.Done():
				if err := b.logout(); err != nil {
					b.notifyError(fmt.Errorf("logout failed: %w", err))
				}
				return
			}
		}
	}()
}

func (b *BetfairClient) reconnect() error {
	token, err := b.login()
	if err != nil {
		return err
	}

	b.mu.Lock()
	b.sessionToken = token
	b.mu.Unlock()

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
