package htmxTasks

import (
	"strconv"
	"todoer/tasks"
	"todoer/templates"
	"net/http"
)

func checkTask(idStr string) tasks.Task {
	id, err := strconv.Atoi(idStr)
	/* User sent stoopid */
	if err != nil {
		panic(err)
	}
	task, err := tasks.Get(id)
	/* No such task */
	if err != nil {
		panic(err)
	}
	return task
}

func GetAll(writer http.ResponseWriter, req *http.Request) {
	taskList := tasks.GetAll()
	data := struct {
		Tasks []tasks.Task
	}{
		Tasks: taskList,
	}
	templates.ExecutePartial(writer, "taskTable", data)
}

func Edit(writer http.ResponseWriter, req *http.Request) {
	task := checkTask(req.PathValue("id"))
	templates.ExecutePartial(writer, "editTaskForm", task)
}

func Clone(writer http.ResponseWriter, req *http.Request) {
	task := checkTask(req.PathValue("id"))
	templates.ExecutePartial(writer, "cloneTaskForm", task)
}
