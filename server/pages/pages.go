package pages

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
	"todoer/tasks"
	"todoer/utils"
)

var (
	layoutsFolder = filepath.Join("server", "templates", "layouts")
	pagesFolder   = filepath.Join("server", "templates", "pages")
	partialsGlob  = filepath.Join("server", "templates", "partials", "*.html")
	modalsGlob    = filepath.Join("server", "templates", "partials", "modals", "*.html")
	pages         map[string]*template.Template
	partials      *template.Template
)

func init() {
	pages = make(map[string]*template.Template)
	partials = template.New("partials").Funcs(TemplateFuncMap)
	template.Must(partials.ParseGlob(partialsGlob))
	template.Must(partials.ParseGlob(modalsGlob))
}

/* FuncMap for HTML templates */
var TemplateFuncMap = template.FuncMap{
	/* Human-readable date formatting */
	"formatDatetime": func(t time.Time) string {
		return t.Format(utils.NiceLookingDatetimeFormat)
	},
	/* Preset dates */
	"getFirstDayOfMonth": func() string {
		fromDate, _ := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
		return fromDate.Format(utils.HTMLDateFormat)
	},
	"getLastDayOfMonth": func() string {
		_, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
		return toDate.Format(utils.HTMLDateFormat)
	},
	"getFirstDayOfPreviousMonth": func() string {
		fromDate, _ := utils.GetMonthBounds(time.Now().Year(), time.Now().AddDate(0, -1, 0).Month())
		return fromDate.Format(utils.HTMLDateFormat)
	},
	"getLastDayOfPreviousMonth": func() string {
		_, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().AddDate(0, -1, 0).Month())
		return toDate.Format(utils.HTMLDateFormat)
	},
	"getToday": func() string {
		now := time.Now()
		return now.Format(utils.HTMLDateFormat)
	},
	"getYesterday": func() string {
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		return yesterday.Format(utils.HTMLDateFormat)
	},
	"parseSortableField": utils.ParseSortableField,
	"parseTaskStatus":    tasks.ParseStatus,
}

func Add(page string, layout string) {
	layoutFullFilename := filepath.Join(layoutsFolder, layout+".html")
	pageFullFilename := filepath.Join(pagesFolder, page+".html")
	/* Create new template, add custom functions */
	newPage := template.New(layout + ".html").Funcs(TemplateFuncMap)
	template.Must(newPage.ParseFiles(layoutFullFilename, pageFullFilename))
	/* Mix in partials */
	template.Must(newPage.ParseGlob(partialsGlob))
	template.Must(newPage.ParseGlob(modalsGlob))
	/* Done */
	pages[page] = newPage
	log.Printf("Added page \"%s/%s\" \n", page, layout)
}

func Execute(writer http.ResponseWriter, name string, data any) {
	/* Check if page exists */
	if page, found := pages[name]; found {
		err := page.Execute(writer, data)
		if err != nil {
			panic(err)
		}
		writer.Header().Set("Content-Type", "text/html")
	} else {
		panic("Invalid page: " + name)
	}
}

func ExecutePartial(writer http.ResponseWriter, name string, data any) {
	err := partials.ExecuteTemplate(writer, name, data)
	if err != nil {
		panic(err)
	}
	writer.Header().Set("Content-Type", "text/html")
}
