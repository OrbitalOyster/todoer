package server

import (
	"log"
	"net/http"
	"todoer/config"
	"todoer/middleware"

	"github.com/klauspost/compress/gzhttp"
)

type RouterEntry func(http.ResponseWriter, *http.Request)
type RouterMap map[string]RouterEntry

func faviconHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/favicon.ico")
}

func Start(routerMap RouterMap) {
	mux := http.NewServeMux()
	/* Static files */
	cssHandler := http.FileServer(http.Dir("static/css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", cssHandler))
	jsHandler := http.FileServer(http.Dir("static/js"))
	mux.Handle("GET /js/", http.StripPrefix("/js/", jsHandler))
	vendorHandler := http.FileServer(http.Dir("static/vendor"))
	mux.Handle("GET /vendor/", http.StripPrefix("/vendor/", vendorHandler))
	/* Favicon */
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	/* Routes */
	for pattern := range routerMap {
		mux.HandleFunc(pattern, routerMap[pattern])
	}
	/* Middleware TODO: Looks like arse */
	middlewared := middleware.Logger(
		middleware.Recovery(
			middleware.Token(
				middleware.Throttle(
					gzhttp.GzipHandler(mux),
				),
			),
		),
	)
	/* Start */
	log.Printf("Starting server on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, middlewared); err != nil {
		panic(err)
	}
}
