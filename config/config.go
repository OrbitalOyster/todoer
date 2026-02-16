package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	Port           string
	CookieName     string
	CookieLifetime int
	JWTSecret      []byte
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Unable to load .env file:", err)
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		log.Panic("Missing PORT variable")
	}

	CookieName = os.Getenv("COOKIE_NAME")
	if CookieName == "" {
		log.Panic("Missing COOKIE_NAME variable")
	}

	CookieLifetime, err = strconv.Atoi(os.Getenv("COOKIE_LIFETIME"))
	if err != nil {
		log.Panic("Error parsing COOKIE_LIFETIME: ", err)
	}

	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(JWTSecret) == 0 {
		log.Panic("Missing JWT_SECRET variable")
	}
}
