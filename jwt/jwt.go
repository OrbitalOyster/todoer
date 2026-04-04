package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"todoer/config"
	"todoer/cookies"
)

type Payload struct {
	UserID     string `json:"user_id"`
	RememberMe bool   `json:"remember_me"`
	PageSize   int    `json:"page_size"`
}

type Claims struct {
	Payload
	jwt.RegisteredClaims
}

func Create(userID string, rememberMe bool, pageSize int) string {
	expirationTime := time.Now()
	if rememberMe {
		expirationTime = expirationTime.Add(
			time.Duration(config.CookieLifetime) * time.Second,
		)
	} else {
		expirationTime = expirationTime.Add(
			time.Duration(config.CookieShortLifetime) * time.Second,
		)
	}
	claims := Claims{
		Payload: Payload{
			UserID:     userID,
			RememberMe: rememberMe,
			PageSize:   pageSize,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	/* Something went terribly wrong */
	if err != nil {
		panic(err)
	}
	return tokenString
}

func Validate(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	return claims, nil
}

func Get(req *http.Request) *Claims {
	cookie := cookies.Get(req)
	/* Should not happen */
	if cookie == "" {
		panic("Empty cookie")
	}
	claims, err := Validate(cookie)
	/* Should not happen */
	if err != nil {
		panic("Empty cookie")
	}
	return claims
}
