package server

import (
	"log"
	"net/http"
	"todoer/config"
	"todoer/routes"
	"todoer/server/middleware"

	"github.com/klauspost/compress/gzhttp"
)

func faviconHandler(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "static/favicon.ico")
}

func Start(routerMap routes.RouterMap) {
	mux := http.NewServeMux()
	/* Static files */
	cssHandler := http.FileServer(http.Dir("static/css"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", middleware.Cache(cssHandler)))
	jsHandler := http.FileServer(http.Dir("static/js"))
	mux.Handle("GET /js/", http.StripPrefix("/js/", middleware.Cache(jsHandler)))
	vendorHandler := http.FileServer(http.Dir("static/vendor"))
	mux.Handle("GET /vendor/", http.StripPrefix("/vendor/", middleware.Cache(vendorHandler)))
	/* Favicon */
	mux.HandleFunc("GET /favicon.ico", faviconHandler)
	/* Routes */
	for pattern := range routerMap {
		mux.HandleFunc(pattern, routerMap[pattern])
	}
	/* Middleware TODO: Looks like arse */
	middlewared := middleware.Logger(
		middleware.Recovery(
			middleware.Auth(
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
