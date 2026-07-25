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
	Request    *http.Request
	Writer     *http.ResponseWriter
	CookieName string
	Secret     []byte
	Claims     claims[T]
}

func Init[T any](
	req *http.Request,
	writer *http.ResponseWriter,
	cookieName string,
	secret []byte,
) Token[T] {
	return Token[T]{
		Request:    req,
		Writer:     writer,
		CookieName: cookieName,
		Secret:     secret,
	}
}

func (token *Token[T]) SetPayload(payload T) {
	token.Claims.Payload = payload
	token.Save()
}

func (token Token[T]) GetPayload() T {
	return token.Claims.Payload
}

func (token *Token[T]) SetLifetime(lifetime int) {
	issuedAt := time.Now()
	expires := issuedAt.Add(time.Duration(lifetime) * time.Second)
	token.Claims = claims[T]{
		Payload: token.Claims.Payload,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
}

func (token Token[T]) GetLifetime() int {
	return int(token.Claims.ExpiresAt.Time.Sub(token.Claims.IssuedAt.Time).Seconds())
}

func (token Token[T]) Save() {
	/* Update lifetime */
	token.SetLifetime(token.GetLifetime())
	/* Create actual jwt */
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token.Claims)
	jwtTokenStr, err := jwtToken.SignedString(token.Secret)
	/* Major screwup */
	if err != nil {
		panic(err)
	}
	/* Done */
	setCookie(
		token.CookieName,
		jwtTokenStr,
		token.Claims.ExpiresAt.Time,
		*token.Writer,
	)
}

func (token *Token[T]) Load() (T, error) {
	var emptyResult T
	cookie := getCookie(token.CookieName, token.Request)
	/* Should not happen */
	if cookie == "" {
		return emptyResult, fmt.Errorf("Empty cookie")
	}
	_, err := jwt.ParseWithClaims(
		cookie,
		&token.Claims,
		func(jwtToken *jwt.Token) (any, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", jwtToken.Header["alg"])
			}
			return token.Secret, nil
		},
	)
	if err != nil {
		return emptyResult, fmt.Errorf("Unable to parse token: %w", err)
	}
	/* Done */
	return token.Claims.Payload, nil
}

func (token *Token[T]) Clear() {
	clearCookie(token.CookieName, *token.Writer)
}
