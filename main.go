package main

import (
	"log"
	"net/http"
	"todoer/api"
	"todoer/config"
	"todoer/middleware"
)

func faviconHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/favicon.ico")
}

func loginHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/html/login.html")
}

func defaultHandler(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.ServeFile(writer, req, "static/html/index.html")
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Nothing here"))
	}
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

	mux.HandleFunc("GET /", defaultHandler)
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	mux.HandleFunc("GET /login", loginHandler)
	mux.HandleFunc("POST /login", api.LoginAttemptHandler)

	cssHandler := http.FileServer(http.Dir("static/css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", cssHandler))

	jsHandler := http.FileServer(http.Dir("static/js"))
	mux.Handle("GET /js/", http.StripPrefix("/js/", jsHandler))

	logged := middlware.Logger(mux)
	authed := middlware.Auth(logged)

	log.Printf("Starting server on port %s", config.Port)
	err := http.ListenAndServe(":"+config.Port, authed)
	if err != nil {
		log.Panic("Unable to start server: ", err)
	}

}
