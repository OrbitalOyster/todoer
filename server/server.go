package server

import (
	"log"
	"net/http"
	"path/filepath"
	"todoer/config"
	"todoer/server/middleware"
	"todoer/server/routes"

	"github.com/klauspost/compress/gzhttp"
)

var (
	cssFolder    = filepath.Join("server", "static", "css")
	jsFolder     = filepath.Join("server", "static", "js")
	vendorFolder = filepath.Join("server", "static", "vendor")
	faviconPath  = filepath.Join("server", "static", "favicon.ico")
)

func Start(routerMap routes.RouterMap) {
	mux := http.NewServeMux()
	/* Static files */
	cssHandler := http.FileServer(http.Dir(cssFolder))
	mux.Handle("GET /css/", http.StripPrefix("/css/", middleware.Cache(cssHandler)))
	jsHandler := http.FileServer(http.Dir(jsFolder))
	mux.Handle("GET /js/", http.StripPrefix("/js/", middleware.Cache(jsHandler)))
	vendorHandler := http.FileServer(http.Dir(vendorFolder))
	mux.Handle("GET /vendor/", http.StripPrefix("/vendor/", middleware.Cache(vendorHandler)))
	/* Favicon */
	mux.HandleFunc("GET /favicon.ico", func(writer http.ResponseWriter, req *http.Request) {
		http.ServeFile(writer, req, faviconPath)
	})
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
