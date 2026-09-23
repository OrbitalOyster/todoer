package token

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/* Actual claims that go into jwt */
type claims[T any] struct {
	Payload T `json:"payload"`
	jwt.RegisteredClaims
}

type Token[T any] struct {
	request    *http.Request
	writer     *http.ResponseWriter
	cookieName string
	secret     []byte
	claims     claims[T]
}

func Init[T any](
	req *http.Request,
	writer *http.ResponseWriter,
	cookieName string,
	secret []byte,
) Token[T] {
	return Token[T]{
		request:    req,
		writer:     writer,
		cookieName: cookieName,
		secret:     secret,
	}
}

func (token *Token[T]) SetPayload(payload T) {
	token.claims.Payload = payload
	token.Save()
}

func (token Token[T]) GetPayload() T {
	return token.claims.Payload
}

func (token *Token[T]) SetLifetime(lifetime int) {
	issuedAt := time.Now()
	expires := issuedAt.Add(time.Duration(lifetime) * time.Second)
	token.claims = claims[T]{
		Payload:   token.claims.Payload,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expires),
	}
}

func (token Token[T]) GetLifetime() int {
	seconds := token.claims.ExpiresAt.Time. /* Get expiration date */
						Sub(token.claims.IssuedAt.Time). /* Substract "issued at" date */
						Seconds()                        /* Convert to seconds */
	return int(seconds) /* Seconds are in float64 */
}

func (token Token[T]) Save() {
	/* Update lifetime */
	token.SetLifetime(token.GetLifetime())
	/* Create actual jwt */
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token.claims)
	jwtTokenStr, err := jwtToken.SignedString(token.secret)
	/* Major screwup */
	if err != nil {
		panic(err)
	}
	/* Done */
	setCookie(
		token.cookieName,
		jwtTokenStr,
		token.claims.ExpiresAt.Time,
		*token.writer,
	)
}

func (token *Token[T]) Load() (result T, err error) {
	cookie := getCookie(token.cookieName, token.request)
	/* Should not happen */
	if cookie == "" {
		return result, fmt.Errorf("Empty cookie")
	}
	_, err = jwt.ParseWithClaims(
		cookie,
		&token.claims,
		func(jwtToken *jwt.Token) (any, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", jwtToken.Header["alg"])
			}
			return token.secret, nil
		},
	)
	if err != nil {
		return result, fmt.Errorf("Unable to parse token: %w", err)
	}
	/* Done */
	return token.claims.Payload, nil
}

func (token *Token[T]) Clear() {
	clearCookie(token.cookieName, *token.writer)
}
