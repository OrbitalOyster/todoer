package middlware

import (
	"net/http"
	"strings"
	"time"
)

func Throttle(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		/* Public routes */
		if strings.HasPrefix(req.URL.Path, "/api/") {
			time.Sleep(time.Second)
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
