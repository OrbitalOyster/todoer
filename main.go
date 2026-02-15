package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	port = 8080
)

func faviconHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/favicon.ico")
}

func loginHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/html/login.html")
}

func loginAttemptHandler(writer http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")

	if username == "orbital" && password == "qwerty" {
		token := createToken("orbital")
		setCookie(writer, token)
		writer.Write([]byte("Yay!"))
	} else {
		writer.Write([]byte("Try again"))
	}
}

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		log.Println("No cookie, redirecting to login")
		http.Redirect(writer, req, "/login", http.StatusSeeOther)
		return
	}
	token, err := validateToken(cookie.Value)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.UserID)
	http.ServeFile(writer, req, "static/html/index.html")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	mux.HandleFunc("GET /login", loginHandler)
	mux.HandleFunc("POST /login", loginAttemptHandler)

	cssHandler := http.FileServer(http.Dir("static/css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", cssHandler))

	jsHandler := http.FileServer(http.Dir("static/js"))
	mux.Handle("GET /js/", http.StripPrefix("/js/", jsHandler))

	loggedMux := loggingMiddleware(mux)

	log.Printf("Starting server on port %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), loggedMux)
	if err != nil {
		panic(err)
	}
}
