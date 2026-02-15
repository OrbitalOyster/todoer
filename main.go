package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	port               = 8080
	cookieName         = "auth"
	jwtSecret          = "JWT_SECRET"
	jwtLifetimeMinutes = 60 * 24 // One day
)

var (
	userName = "orbital"
	password = "qwerty"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateToken(userID string) string {
	expirationTime := time.Now().Add(time.Duration(jwtLifetimeMinutes) * time.Minute)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func SetCookie(writer http.ResponseWriter, value string) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    value,
		Expires:  time.Now().Add(time.Duration(jwtLifetimeMinutes) * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(writer, &cookie)
}

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		log.Println("No cookie")
		return
	}
	token, err := validateToken(cookie.Value)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.UserID)
	http.ServeFile(writer, req, "static/html/index.html")
}

func faviconHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/favicon.ico")
}
func loginHandler(writer http.ResponseWriter, req *http.Request) {
	userID := "User"
	token := CreateToken(userID)
	SetCookie(writer, token)
	http.ServeFile(writer, req, "static/html/login.html")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/login", loginHandler)

	cssHandler := http.FileServer(http.Dir("static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))

	log.Printf("Starting server on port %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}
