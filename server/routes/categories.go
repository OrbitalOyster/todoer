package routes

import (
	"net/http"
	"todoer/server/pages"
	"todoer/utils"
)

func GetCategoriesPage(writer http.ResponseWriter, req *http.Request) {
	payload := utils.GetTokenPayload(req)
	pages.Execute(
		writer,
		"categories",
		struct {
			Title   string
			Payload utils.Payload
		}{
			Title:   "todoer - Categories",
			Payload: payload,
		},
	)
}
