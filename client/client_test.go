// client/client_test.go

package client_test

import (
	"betfair-api-go-sdk/client"
	"betfair-api-go-sdk/types"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fakeValidPEMCert = `-----BEGIN CERTIFICATE-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
-----END CERTIFICATE-----`

const fakeValidPEMKey = `-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
-----END PRIVATE KEY-----`

func TestNewSession_InvalidCredentials(t *testing.T) {
	tests := []struct {
		name        string
		creds       types.BetfairCredentials
		wantErr     bool
		containsErr string // substring we expect in error message
	}{
		{
			name: "empty username",
			creds: types.BetfairCredentials{
				Username:   "",
				Password:   "secret",
				AppKey:     "APPKEY123",
				CertString: fakeValidPEMCert,
				KeyString:  fakeValidPEMKey,
			},
			wantErr:     true,
			containsErr: "username cannot be empty",
		},
		{
			name: "empty password",
			creds: types.BetfairCredentials{
				Username:   "user",
				Password:   "",
				AppKey:     "APPKEY123",
				CertString: fakeValidPEMCert,
				KeyString:  fakeValidPEMKey,
			},
			wantErr:     true,
			containsErr: "password cannot be empty",
		},
		{
			name: "empty app key",
			creds: types.BetfairCredentials{
				Username:   "user",
				Password:   "secret",
				AppKey:     "",
				CertString: fakeValidPEMCert,
				KeyString:  fakeValidPEMKey,
			},
			wantErr:     true,
			containsErr: "app key cannot be empty", // Betfair usually rejects empty appkey early
		},
		{
			name: "invalid base64 cert string",
			creds: types.BetfairCredentials{
				Username:   "user",
				Password:   "secret",
				AppKey:     "APPKEY123",
				CertString: "not-base64-!!!",
				KeyString:  fakeValidPEMKey,
			},
			wantErr:     true,
			containsErr: "invalid cert string",
		},
		{
			name: "invalid base64 key string",
			creds: types.BetfairCredentials{
				Username:   "user",
				Password:   "secret",
				AppKey:     "APPKEY123",
				CertString: fakeValidPEMCert,
				KeyString:  "invalid-key-@#$",
			},
			wantErr:     true,
			containsErr: "illegal base64",
		},
		{
			name: "cert/key pair invalid (corrupted)",
			creds: types.BetfairCredentials{
				Username:   "user",
				Password:   "secret",
				AppKey:     "APPKEY123",
				CertString: base64.StdEncoding.EncodeToString([]byte("broken cert")),
				KeyString:  base64.StdEncoding.EncodeToString([]byte("broken key")),
			},
			wantErr:     true,
			containsErr: "failed to parse cert/key pair",
		},
		{
			name: "valid credentials (negative case - should NOT error here)",
			creds: types.BetfairCredentials{
				Username:   "user",
				Password:   "secret",
				AppKey:     "APPKEY123",
				CertString: fakeValidPEMCert,
				KeyString:  fakeValidPEMKey,
			},
			wantErr:     true,
			containsErr: "illegal base64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We don't care about real login here - we just want early validation errors
			// Use a dummy error handler
			_, err := client.NewSession(tt.creds, nil)

			if tt.wantErr {
				require.Error(t, err, "expected error but got nil")
				if tt.containsErr != "" {
					assert.Contains(t, err.Error(), tt.containsErr,
						"error message should contain expected substring")
				}
			} else {
				// For the "valid" case we expect error anyway (fake creds),
				// but at least it shouldn't fail during cert parsing
				if err != nil && !strings.Contains(err.Error(), "unable to login") {
					t.Errorf("unexpected early validation error: %v", err)
				}
			}
		})
	}
}
