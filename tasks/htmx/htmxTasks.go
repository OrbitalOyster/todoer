package htmxTasks

import (
	"log"
	"net/http"
	"strconv"
	"todoer/tasks"
	"todoer/templates"
	"todoer/toasts"
)

const (
	defaultSize = 10
	maxPageSize = 50
)

func Get(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "task", task)
}

func GetAll(writer http.ResponseWriter, req *http.Request) {

	/* Parse REST query */
	query := req.URL.Query()
	filter := query.Get("filter")
	sort := query.Get("sort")
	size, err := strconv.Atoi(query.Get("size"))
	if err != nil {
		size = defaultSize
	}
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}

	log.Printf("sort: %s, page: %d, size: %d, filter: %s", sort, page, size, filter)
	data := struct {
		Tasks []tasks.Task
	}{
		Tasks: tasks.GetAll(filter, size, page),
	}
	templates.ExecutePartial(writer, "taskTable", data)
}

func Edit(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "editTaskForm", task)
}

func Clone(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "cloneTaskForm", task)
}

func Patch(writer http.ResponseWriter, req *http.Request) {
	idStr, description := req.FormValue("id"), req.FormValue("description")
	if description == "bogus" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bogus description"))
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}
	tasks.Update(id, description)
	task := tasks.Check(idStr)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "Task "+idStr, "Success")
	templates.ExecutePartial(writer, "task-oob", task)
}
