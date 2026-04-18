package utils

import (
	"slices"
	"strings"
)

var publicURIs = []string{
	"/login",
	"/favicon.ico",
	"/css/reset.css",
	"/css/style.css",
	"/js/script.js",
}

func IsPublicURL(URL string) bool {
	return slices.Contains(publicURIs, URL) || strings.HasPrefix(URL, "/vendor/")
}
