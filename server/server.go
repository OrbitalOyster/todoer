package server

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

func Start() {
	mux := http.NewServeMux()
	/* Static files */
	cssHandler := http.FileServer(http.Dir("static/css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", cssHandler))
	jsHandler := http.FileServer(http.Dir("static/js"))
	mux.Handle("GET /js/", http.StripPrefix("/js/", jsHandler))
	/* Routes */
	mux.HandleFunc("GET /", defaultHandler)
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	mux.HandleFunc("GET /login", loginHandler)
	mux.HandleFunc("POST /login", api.LoginAttemptHandler)
	/* Middleware */
	middlewared := middlware.Auth(middlware.Logger(mux))
	/* Start */
	log.Printf("Starting server on port %s", config.Port)
	err := http.ListenAndServe(":"+config.Port, middlewared)
	if err != nil {
		panic(err)
	}
}
