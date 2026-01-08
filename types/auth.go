// types/auth.go

package types

type BetfairCredentials struct {
	Username   string
	Password   string
	AppKey     string
	CertString string // Base64 encoded
	KeyString  string // Base64 encoded

	ProxyUrl *string //optional
}
