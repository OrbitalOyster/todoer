package api

import (
	"log"
	"net/http"
	"strconv"
)

func SetTaskTablePageSize(writer http.ResponseWriter, req *http.Request) {
	/* Check if form is ok */
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	size, err := strconv.ParseUint(req.FormValue("size"), 10, 64)
	if err != nil {
		size = 10 /* TODO: Magic number */
	}
	log.Printf("Setting page size to %d", size)
	
	writer.Write([]byte("Hello"))
}
