package routes

import (
	"net/http"
	"todoer/server/pages"
	"todoer/users"
)

func GetUsersPage(writer http.ResponseWriter, req *http.Request) {
	allUsers := users.GetAllUsers()
	pages.Execute(
		writer,
		"users",
		struct {
			Title string
			Users []users.User
		}{
			Title: "Users",
			Users: allUsers,
		},
	)
}
