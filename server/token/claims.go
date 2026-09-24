package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/* Actual claims that go into token */
type claims[T any] struct {
	Payload T `json:"payload"`
	jwt.RegisteredClaims
}

func (claims *claims[T]) SetLifetime(lifetime int) {
	issuedAt := time.Now()
	expires := issuedAt.Add(time.Duration(lifetime) * time.Second)
	claims.IssuedAt = jwt.NewNumericDate(issuedAt)
	claims.NotBefore = claims.IssuedAt
	claims.ExpiresAt = jwt.NewNumericDate(expires)
}

func (claims claims[T]) GetLifetime() (seconds int) {
	secondsFloat := claims.ExpiresAt.Time. /* Get expiration date */
						Sub(claims.IssuedAt.Time). /* Substract "issued at" date */
						Seconds()                  /* Convert to seconds */
	return int(secondsFloat)
}

func (claims *claims[T]) SetPayload(payload T) {
	claims.Payload = payload
}

func (claims claims[T]) GetPayload() T {
	return claims.Payload
}
