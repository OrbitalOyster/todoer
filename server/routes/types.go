package routes

import (
	"net/http"
)

type RouterEntry func(http.ResponseWriter, *http.Request)
type RouterMap map[string]RouterEntry
