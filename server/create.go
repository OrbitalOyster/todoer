package server

import (
	"log"
	"todoer/api"
	"net/http"
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

func Create(port string) {
	mux := http.NewServeMux()
	cssHandler := http.FileServer(http.Dir("static/css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", cssHandler))
	jsHandler := http.FileServer(http.Dir("static/js"))
	mux.Handle("GET /js/", http.StripPrefix("/js/", jsHandler))

	mux.HandleFunc("GET /", defaultHandler)
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	mux.HandleFunc("GET /login", loginHandler)
	mux.HandleFunc("POST /login", api.LoginAttemptHandler)

	a0 := middlware.Logger(mux)
	a1 := middlware.Auth(a0)

	err := http.ListenAndServe(":"+port, a1)
	if err != nil {
		log.Panic("Unable to start server: ", err)
	}
}
