package toasts

import (
	"net/http"
	"time"
	"todoer/templates"
)

const timeFormat = "2.01.2006 15:04:05"

func setToastHeaders(writer http.ResponseWriter) {
	writer.Header().Set("HX-Trigger-After-Settle", "toast")
	writer.Header().Set("HX-Retarget", ".toast-container")
	writer.Header().Set("HX-Reswap", "beforeend")
}

func Success(writer http.ResponseWriter, title string, msg string) {
	setToastHeaders(writer)
	options := map[string]string {
		"BorderClass": "border-success-subtle",
		"Autohide": "true",
		"HeaderColor": "bg-success-subtle",
		"IconColor": "text-success",
		"IconClass": "bi-hand-thumbs-up-fill",
		"Title": title,
		"ProgressBarClass": "bg-success",
		"Content": msg,
	}
	options["Time"] = time.Now().Format(timeFormat)
	templates.ExecutePartial(writer, "toast", options)	
}

func Warning(writer http.ResponseWriter, title string, msg string) {
	setToastHeaders(writer)
	options := map[string]string {
		"BorderClass": "border-warning-subtle",
		"Autohide": "false",
		"HeaderColor": "bg-warning-subtle",
		"IconColor": "text-warning",
		"IconClass": "bi-exclamation-triangle-fill",
		"Title": title,
		"ProgressBarClass": "d-none",
		"Content": msg,
	}
	options["Time"] = time.Now().Format(timeFormat)
	templates.ExecutePartial(writer, "toast", options)	
}
