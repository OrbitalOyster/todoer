package routes

import (
	"net/http"
	"todoer/server/pages"
	"todoer/server/token"
	"todoer/utils"
)

func GetCategoriesPage(writer http.ResponseWriter, req *http.Request) {
	payload := req.Context(). /* Get context from request */
					Value("token").(*token.Token[utils.Payload]). /* Get "token" field */
					GetPayload()
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
