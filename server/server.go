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
	cssFolder   = filepath.Join("static", "css")
	jsFolder    = filepath.Join("static", "js")
	imgFolder   = filepath.Join("static", "img")
	faviconPath = filepath.Join("static", "favicon.ico")

	bootstrapFolder      = filepath.Join("external", "bootstrap", "dist")
	bootstrapIconsFolder = filepath.Join("external", "bootstrap-icons", "font")
	bootswatchFolder     = filepath.Join("external", "bootswatch", "dist")
	htmxFolder           = filepath.Join("external", "htmx", "dist")
)

func Start(routerMap routes.RouterMap) {
	mux := http.NewServeMux()
	/* Static files */
	cssHandler := http.FileServer(http.Dir(cssFolder))
	mux.Handle("GET /css/", http.StripPrefix("/css/", middleware.Cache(cssHandler)))
	jsHandler := http.FileServer(http.Dir(jsFolder))
	mux.Handle("GET /js/", http.StripPrefix("/js/", middleware.Cache(jsHandler)))
	/* External handlers */
	bootstrapHandler := http.FileServer(http.Dir(bootstrapFolder))
	mux.Handle("GET /bootstrap/", http.StripPrefix("/bootstrap/", middleware.Cache(bootstrapHandler)))
	bootstrapIconsHandler := http.FileServer(http.Dir(bootstrapIconsFolder))
	mux.Handle("GET /bootstrap-icons/", http.StripPrefix("/bootstrap-icons/", middleware.Cache(bootstrapIconsHandler)))
	bootswatchHandler := http.FileServer(http.Dir(bootswatchFolder))
	mux.Handle("GET /bootswatch/", http.StripPrefix("/bootswatch/", middleware.Cache(bootswatchHandler)))
	htmxHandler := http.FileServer(http.Dir(htmxFolder))
	mux.Handle("GET /htmx/", http.StripPrefix("/htmx/", middleware.Cache(htmxHandler)))
	/* Images */
	imgHandler := http.FileServer(http.Dir(imgFolder))
	mux.Handle("GET /img/", http.StripPrefix("/img/", middleware.Cache(imgHandler)))
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
				middleware.ParseForm(
					middleware.Throttle(
						gzhttp.GzipHandler(mux),
					),
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
