package middleware

import (
	"log"
	"net/http"
)

func ParseForm(next http.Handler) http.Handler {
	handler := func(writer http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			panic("Haxxor alert!")
		} else {
			if len(req.Form) > 0 {
				/* TODO: Hide passwords */
				log.Printf("[DEBUG] Form data: %#v\n", req.Form)
			}
		}
		next.ServeHTTP(writer, req)
	}
	return http.HandlerFunc(handler)
}
