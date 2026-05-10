// auth/auth.go
package auth

import "time"

type Auth struct {
	Tokens *TokenService
}

func New(secret string, accessDuration, refreshDuration time.Duration) *Auth {
	return &Auth{
		Tokens: NewTokenService(secret, accessDuration, refreshDuration),
	}
}
