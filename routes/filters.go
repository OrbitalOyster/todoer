package routes

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"todoer/config"
	"todoer/cookies"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
)

var pageSizes = []int{5, 10, 25, 50}

func SetTaskTablePageSize(writer http.ResponseWriter, req *http.Request) {
	/* Check if form is ok */
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(req.FormValue("size"))
	if err != nil || !slices.Contains(pageSizes, size) {
		size = config.DefaultPageSize
	}
	log.Printf("Setting page size to %d", size)
	/* Update token/cookies */
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.PageSize = size
	token := jwt.Set(*payload)
	cookies.Set(writer, token, payload.RememberMe)

	/* Return updated task table */
	templates.ExecutePartial(writer, "task-table-body", tasks.GetFromPayload(*payload))
}
