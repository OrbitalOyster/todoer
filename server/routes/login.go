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

		funkyToken := token.CreateFunkyToken[token.Payload](
			req,
			&writer,
			config.CookieName,
			config.JWTSecret,
			config.CookieShortLifetime,
		)
		fromDate, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
		freshPayload := token.Payload{
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
		funkyToken.SetPayload(freshPayload)
		funkyToken.Save()

		// token.CreateFresh(username, rememberMe, writer)
		writer.Header().Set("HX-Redirect", "/")
		log.Printf("User %s logged in", username)
	} else {
		toasts.Info(writer, "Login failed", "Try again")
	}
}

func Logout(writer http.ResponseWriter, req *http.Request) {

	funkyToken := token.CreateFunkyToken[token.Payload](
		req,
		&writer,
		config.CookieName,
		config.JWTSecret,
		config.CookieShortLifetime,
	)
	// user := token.Get(req).UserID
	// token.Clear(writer)

	user := funkyToken.GetPayload().UserID
	funkyToken.Clear()

	writer.Header().Set("HX-Redirect", "/login")
	log.Printf("User %s logged out", user)
}
