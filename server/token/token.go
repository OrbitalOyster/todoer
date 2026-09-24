package token

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Token[T any] struct {
	request    *http.Request
	writer     *http.ResponseWriter
	cookieName string
	secret     []byte
	claims[T]
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

func (token *Token[T]) Validate() (err error) {
	cookie := getCookie(token.cookieName, token.request)
	/* Should not happen */
	if cookie == "" {
		return fmt.Errorf("Empty cookie")
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
		return fmt.Errorf("Unable to parse token: %w", err)
	}
	/* Done */
	return nil
}

func (token *Token[T]) Clear() {
	clearCookie(token.cookieName, *token.writer)
}
