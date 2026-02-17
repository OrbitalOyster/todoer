package routes

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"todoer/jwt"
)

var layout *template.Template

func init()  {
	log.Println("Initializing template...")
	file, err := os.ReadFile("templates/root.html")
	if err != nil {
		panic(err)
	}
	parsed, err := template.New("test").Parse(string(file))
	if err != nil {
		panic(err)
	}
	layout = parsed
}

func Foo(writer http.ResponseWriter, req *http.Request) {
	claims := jwt.Get(req)
	data := struct {
		Title string
		Username string
	} {
		Title: "My template",
		Username: claims.UserID,
	}

	layout.Execute(writer, data)
}
