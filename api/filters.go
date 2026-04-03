package api

import (
	"log"
	"net/http"
)

func SetTaskTablePageSize(writer http.ResponseWriter, req *http.Request) {
	/* Check if form is ok */
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	size := req.FormValue("size")
	log.Printf("Setting page size to %s", size)
	writer.Write([]byte("Hello"))
}
