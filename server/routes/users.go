package routes

import (
	"net/http"
	"todoer/server/pages"
	"todoer/users"
	"todoer/utils"
)

func GetUsersPage(writer http.ResponseWriter, req *http.Request) {
	payload := utils.GetTokenPayload(req)
	allUsers := users.GetAllUsers()
	pages.Execute(
		writer,
		"users",
		struct {
			Title   string
			Users   []users.User
			Payload utils.Payload
		}{
			Title:   "todoer - Users",
			Payload: payload,
			Users:   allUsers,
		},
	)
}
