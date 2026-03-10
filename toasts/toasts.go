package toasts

import (
	"maps"
	"net/http"
	"strconv"
	"time"
	"todoer/templates"
)

type ToastSeverity int
const (
	success ToastSeverity = iota
	info
	warning
	danger
)

var successOpts = map[string]string{
	"BorderClass":      "border-success-subtle",
	"Autohide":         "true",
	"HeaderColor":      "bg-success-subtle",
	"IconColor":        "text-success",
	"IconClass":        "bi-hand-thumbs-up-fill",
	"ProgressBarColor": "bg-success",
}

var infoOpts = map[string]string{
	"BorderClass":      "border-info-subtle",
	"Autohide":         "true",
	"HeaderColor":      "bg-info-subtle",
	"IconColor":        "text-info",
	"IconClass":        "bi-info-circle-fill",
	"ProgressBarColor": "bg-info",
}

var warningOpts = map[string]string{
	"BorderClass":      "border-warning-subtle",
	"Autohide":         "false",
	"HeaderColor":      "bg-warning-subtle",
	"IconColor":        "text-warning",
	"IconClass":        "bi-exclamation-triangle-fill",
	"ProgressBarColor": "bg-warning",
}

var dangerOpts = map[string]string{
	"BorderClass":      "border-danger-subtle",
	"Autohide":         "false",
	"HeaderColor":      "bg-danger-subtle",
	"IconColor":        "text-danger",
	"IconClass":        "bi-x-octagon-fill",
	"ProgressBarColor": "bg-danger",
}

const timeFormat = "2.01.2006 15:04:05"

func execute(writer http.ResponseWriter, severity ToastSeverity, title string, msg string) {
	var options map[string]string
	switch severity {
	case success:
		options = maps.Clone(successOpts)
	case info:
		options = maps.Clone(infoOpts)
	case warning:
		options = maps.Clone(warningOpts)
	case danger:
		options = maps.Clone(dangerOpts)
	default:
		panic("Invalid toast severity: " + strconv.Itoa(int(severity)))
	}
	options["Title"] = title
	options["Time"] = time.Now().Format(timeFormat)
	options["Content"] = msg
	writer.Header().Set("HX-Trigger-After-Settle", "toast")
	writer.Header().Set("HX-Retarget", ".toast-container")
	writer.Header().Set("HX-Reswap", "beforeend")
	templates.ExecutePartial(writer, "toast", options)
}

func Success(writer http.ResponseWriter, title string, msg string) {
	execute(writer, success, title, msg)
}

func Info(writer http.ResponseWriter, title string, msg string) {
	execute(writer, 7, title, msg)
}

func Warning(writer http.ResponseWriter, title string, msg string) {
	execute(writer, warning, title, msg)
}

func Danger(writer http.ResponseWriter, title string, msg string) {
	execute(writer, danger, title, msg)
}
