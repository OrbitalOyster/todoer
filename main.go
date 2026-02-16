package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"todoer/api"
	"todoer/server"
)

var (
	port       string
	cookieName string
)

func faviconHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/favicon.ico")
}

func loginHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/html/login.html")
}

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	cookie := server.GetCookie(req)
	if cookie == "" {
		http.Redirect(writer, req, "/login", http.StatusSeeOther)
		return
	}
	token, err := server.ValidateToken(cookie)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.UserID)
	http.ServeFile(writer, req, "static/html/index.html")
}

func main() {
	/* Error handler */
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Println("Crashed:", recovered)
		}
		log.Println("Bye")
	}()

	/* Load .env file */
	dotEnvErr := godotenv.Load()
	if dotEnvErr != nil {
		log.Panic("Missing .env file")
	}

	port = os.Getenv("PORT")
	if port == "" {
		log.Panic("Missing PORT variable")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	mux.HandleFunc("GET /login", loginHandler)
	mux.HandleFunc("POST /login", api.LoginAttemptHandler)

	cssHandler := http.FileServer(http.Dir("static/css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", cssHandler))

	jsHandler := http.FileServer(http.Dir("static/js"))
	mux.Handle("GET /js/", http.StripPrefix("/js/", jsHandler))

	loggedMux := loggingMiddleware(mux)

	log.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(":"+port, loggedMux)
	if err != nil {
		log.Panic("Unable to start server: ", err)
	}

}
