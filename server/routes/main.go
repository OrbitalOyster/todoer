package routes

import (
	"net/http"
	"todoer/server/pages"
	"todoer/server/token"
	"todoer/utils"
)

func GetMainPage(writer http.ResponseWriter, req *http.Request) {
	payload := req.Context(). /* Get context from request */
					Value("token").(*token.Token[utils.Payload]). /* Get "token" field */
					GetPayload()                                  /* Load actual payload */
	pages.Execute(writer, "main", struct {
		Title   string
		Payload utils.Payload
	}{
		Title:   "todoer",
		Payload: payload,
	})
}
