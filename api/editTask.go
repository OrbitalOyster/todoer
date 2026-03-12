package api

import (
	"net/http"
)

func EditTask(writer http.ResponseWriter, req *http.Request)  {
	writer.Write([]byte("Hello"))	
}
