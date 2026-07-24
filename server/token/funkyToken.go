package token

import (
	"fmt"
	"net/http"
	"time"
	"todoer/config"

	"github.com/golang-jwt/jwt/v5"
)

type FunkyClaims[T any] struct {
	Payload T `json:"payload"`
	jwt.RegisteredClaims
}

type FunkyToken[T any] struct {
	Request    *http.Request
	Writer     *http.ResponseWriter
	CookieName string
	Secret     []byte
	Lifetime   int
	Claims     FunkyClaims[T]
}

func CreateFunkyToken[T any](
	req *http.Request,
	writer *http.ResponseWriter,
	cookieName string,
	secret []byte,
	lifetime int) FunkyToken[T] {
	return FunkyToken[T]{
		Request:    req,
		Writer:     writer,
		CookieName: cookieName,
		Secret:     secret,
		Lifetime:   lifetime,
	}
}

func (funkyToken *FunkyToken[T]) SetPayload(payload T) {
	funkyToken.Claims.Payload = payload
}

func (funkyToken FunkyToken[T]) GetPayload() T {
	return funkyToken.Claims.Payload
}

func (funkyToken FunkyToken[T]) Save() {
	expires := time.Now().Add(time.Duration(funkyToken.Lifetime) * time.Second)
	funkyToken.Claims = FunkyClaims[T]{
		Payload: funkyToken.Claims.Payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, funkyToken.Claims)
	tokenStr, err := token.SignedString(config.JWTSecret)
	/* Major screwup */
	if err != nil {
		panic(err)
	}
	setCookie(funkyToken.CookieName, tokenStr, expires, *funkyToken.Writer)
}

func (funkyToken *FunkyToken[T]) Load() error {
	cookie := getCookie(funkyToken.CookieName, funkyToken.Request)
	/* Should not happen */
	if cookie == "" {
		return fmt.Errorf("Empty cookie")
	}
	_, err := jwt.ParseWithClaims(
		cookie,
		&funkyToken.Claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return funkyToken.Secret, nil
		},
	)
	if err != nil {
		return fmt.Errorf("Unable to parse funky token: %w", err)
	}
	fmt.Printf("debug: %v\n", funkyToken.Claims.Payload)
	return nil
}

func (funkyToken *FunkyToken[T]) Clear() {
	clearCookie(funkyToken.CookieName, *funkyToken.Writer)
}
