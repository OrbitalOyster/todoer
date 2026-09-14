package routes

import (
	"net/http"
	"todoer/server/pages"
	"todoer/utils"
)

func GetMainPage(writer http.ResponseWriter, req *http.Request) {
	payload := utils.GetTokenPayload(req)
	pages.Execute(writer, "main", struct {
		Title   string
		Payload utils.Payload
	}{
		Title:   "todoer",
		Payload: payload,
	})
}
