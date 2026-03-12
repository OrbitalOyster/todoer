package api

import (
	"net/http"
	"todoer/templates"
)

func EditTask(writer http.ResponseWriter, req *http.Request)  {
	templates.ExecutePartial(writer, "editTaskModalHTMX", nil)
}
