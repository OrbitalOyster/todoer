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
	Lifetime   int
	Claims     claims[T]
}

func Init[T any](
	req *http.Request,
	writer *http.ResponseWriter,
	cookieName string,
	secret []byte,
	lifetime int) Token[T] {
	return Token[T]{
		Request:    req,
		Writer:     writer,
		CookieName: cookieName,
		Secret:     secret,
		Lifetime:   lifetime,
	}
}

func (token *Token[T]) SetPayload(payload T) {
	token.Claims.Payload = payload
}

func (token Token[T]) GetPayload() T {
	return token.Claims.Payload
}

func (token Token[T]) Save() {
	expires := time.Now().Add(time.Duration(token.Lifetime) * time.Second)
	token.Claims = claims[T]{
		Payload: token.Claims.Payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token.Claims)
	jwtTokenStr, err := jwtToken.SignedString(token.Secret)
	/* Major screwup */
	if err != nil {
		panic(err)
	}
	setCookie(token.CookieName, jwtTokenStr, expires, *token.Writer)
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
	return token.Claims.Payload, nil
}

func (token *Token[T]) Clear() {
	clearCookie(token.CookieName, *token.Writer)
}
