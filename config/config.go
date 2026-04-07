package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Port                string
	CookieName          string
	CookieLifetime      int
	CookieShortLifetime int
	JWTSecret           []byte
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Unable to load .env file: ", err)
	}
	/* HTTP port */
	Port = os.Getenv("PORT")
	if Port == "" {
		log.Panic("Missing PORT variable")
	}
	/* Cookie name */
	CookieName = os.Getenv("COOKIE_NAME")
	if CookieName == "" {
		log.Panic("Missing COOKIE_NAME variable")
	}
	/* Cookie lifetime with "remember_me" option */
	CookieLifetime, err = strconv.Atoi(os.Getenv("COOKIE_LIFETIME"))
	if err != nil {
		log.Panic("Error parsing COOKIE_LIFETIME: ", err)
	}
	/* Cookie lifetime without "remember_me" option */
	CookieShortLifetime, err = strconv.Atoi(os.Getenv("COOKIE_SHORT_LIFETIME"))
	if err != nil {
		log.Panic("Error parsing COOKIE_SHORT_LIFETIME: ", err)
	}
	/* JWT secret */
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(JWTSecret) == 0 {
		log.Panic("Missing JWT_SECRET variable")
	}
}
