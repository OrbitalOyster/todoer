package routes

import (
	"log"
	"net/http"
	"time"
	"todoer/config"
	"todoer/server/pages"
	"todoer/server/toasts"
	"todoer/server/token"
	"todoer/utils"
)

/* GET */
func GetLoginPage(writer http.ResponseWriter, req *http.Request) {
	data := struct{ Title string }{"Login"}
	pages.Execute(writer, "login", data)
}

/* POST */
func LoginAttempt(writer http.ResponseWriter, req *http.Request) {
	/* Credentials mock up */
	username, password, rememberMeStr, rememberMe :=
		req.FormValue("username"),
		req.FormValue("password"),
		req.FormValue("remember-me"),
		false
	if rememberMeStr == "true" {
		rememberMe = true
	}
	/* Auth mockup */
	if username == "admin" && password == "password" {
		lifetime := config.CookieShortLifetime
		if rememberMe {
			lifetime = config.CookieLifetime
		}
		token := token.Init[utils.Payload](
			req,
			&writer,
			config.CookieName,
			config.JWTSecret,
			lifetime,
		)
		fromDate, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
		freshPayload := utils.Payload{
			UserID:     username,
			RememberMe: rememberMe,
			PageSize:   config.DefaultPageSize,
			Page:       1,
			SearchBy:   "",
			SortBy:     utils.Datetime,
			SortAsc:    true,
			FromDate:   fromDate.Format(utils.HTMLDateFormat),
			ToDate:     toDate.Format(utils.HTMLDateFormat),
		}
		token.SetPayload(freshPayload)
		token.Save()
		writer.Header().Set("HX-Redirect", "/")
		log.Printf("User %s logged in", username)
	} else {
		toasts.Info(writer, "Login failed", "Try again")
	}
}

func Logout(writer http.ResponseWriter, req *http.Request) {
	token := req.Context(). /* Get context from request */
				Value("token").(*token.Token[utils.Payload]) /* Get "token" field */
	payload := token.GetPayload() /* Load actual payload */
	token.Clear()
	writer.Header().Set("HX-Redirect", "/login")
	log.Printf("User %s logged out", payload.UserID)
}
