// types/auth.go

package types

import "net/url"

type BetfairCredentials struct {
	Username   string
	Password   string
	AppKey     string
	CertString string   // Base64 encoded
	KeyString  string   // Base64 encoded
	ProxyUrl   *url.URL //optional
}
