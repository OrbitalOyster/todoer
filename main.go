package main

import (
	"fmt"
	"log"
	"net/http"
	"todoer/api"
	"todoer/config"
	"todoer/server"
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
		http.Redirect(writer, req, "/login", http.StatusSeeOther)
		return
		// panic(err)
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

	config.Load()

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

	log.Printf("Starting server on port %s", config.Port)
	err := http.ListenAndServe(":"+config.Port, loggedMux)
	if err != nil {
		log.Panic("Unable to start server: ", err)
	}

}
