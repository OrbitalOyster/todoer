package api

import (
	"log"
	"net/http"
	"todoer/cookies"
	"todoer/jwt"
)

func Logout(writer http.ResponseWriter, req *http.Request) {
	claims, err := jwt.Validate(cookies.Get(req))
	if err != nil {
		panic(err)
	}
	cookies.Clear(writer)
	writer.Header().Set("HX-Redirect", "/login")
	log.Printf("User %s logged out", claims.UserID)
}
