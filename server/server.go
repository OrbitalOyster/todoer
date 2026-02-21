package server

import (
	"log"
	"net/http"
	"todoer/config"
	"todoer/middleware"
)

type RouterEntry func(http.ResponseWriter, *http.Request)
type RouterMap map[string] RouterEntry

func faviconHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/favicon.ico")
}

func Start(routerMap RouterMap) {
	mux := http.NewServeMux()
	/* Static files */
	cssHandler := http.FileServer(http.Dir("static/css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", cssHandler))
	/* Favicon */
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	/* Routes */
	for pattern := range routerMap {
		mux.HandleFunc(pattern, routerMap[pattern])
	}
	/* Middleware */
	middlewared := middlware.Auth(middlware.Logger(mux))
	/* Start */
	log.Printf("Starting server on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, middlewared); err != nil {
		panic(err)
	}
}
