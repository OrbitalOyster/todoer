package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"todoer/tasks"
	"todoer/utils"
)

type TasksQuery[T tasks.TaskFieldName] struct {
	Page     int
	Size     int
	SearchBy string
	FromDate string
	ToDate   string
	SortBy   T
	SortDesc bool
}

func defaultTaskQuery() TasksQuery[tasks.TaskFieldName] {
	fromDate, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	return TasksQuery[tasks.TaskFieldName]{
		Page:     1,
		Size:     defaultPageSize,
		SearchBy: "",
		FromDate: fromDate.Format(utils.HTMLDateFormat),
		ToDate:   toDate.Format(utils.HTMLDateFormat),
		SortBy:   tasks.Datetime,
		SortDesc: false,
	}
}

/*
func (taskQuery *TasksQuery[T]) Default() {
	fromDate, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	taskQuery.Page = 1
	taskQuery.Size = defaultPageSize
	taskQuery.SearchBy = ""
	taskQuery.FromDate = fromDate.Format(utils.HTMLDateFormat)
	taskQuery.ToDate = toDate.Format(utils.HTMLDateFormat)
	taskQuery.SortBy = T(tasks.Datetime)
	taskQuery.SortDesc = false
}
*/

func CreateQueryFromRequest(req *http.Request) TasksQuery[tasks.TaskFieldName] {
	query := defaultTaskQuery()
	/* HTMX request */
	if hxCurrentUrl := req.Header.Get("HX-Current-URL"); hxCurrentUrl != "" {
		/* Current browser query */
		url, err := url.Parse(hxCurrentUrl)
		if err != nil {
			panic(err)
		}
		query.Parse(url.RawQuery)
	}
	/* Actual request */
	query.Parse(req.URL.RawQuery)
	return query
}

func (taskQuery *TasksQuery[T]) Parse(rawQuery string) {
	parsed, err := url.ParseQuery(rawQuery)
	if err != nil {
		return
	}
	if parsed.Has("page") {
		if parsedPage, err := strconv.Atoi(parsed.Get("page")); err == nil {
			taskQuery.Page = parsedPage
		}
	}
	if parsed.Has("size") {
		if parsedSize, err := strconv.Atoi(parsed.Get("size")); err == nil {
			taskQuery.Size = parsedSize
		}
	}
	if parsed.Has("searchBy") {
		taskQuery.SearchBy = parsed.Get("searchBy")
	}

	defaultFromDate, defaultToDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	if parsed.Has("from") {
		if _, err := time.Parse(utils.HTMLDateFormat, parsed.Get("from")); err != nil {
			taskQuery.FromDate = defaultFromDate.Format(utils.HTMLDateFormat)
		} else {
			taskQuery.FromDate = parsed.Get("from")
		}
	}
	if parsed.Has("to") {
		if _, err := time.Parse(utils.HTMLDateFormat, parsed.Get("to")); err != nil {
			taskQuery.ToDate = defaultToDate.Format(utils.HTMLDateFormat)
		} else {
			taskQuery.ToDate = parsed.Get("to")
		}
	}
	/* Sorting */
	if parsed.Has("sortBy") {
		taskQuery.SortBy = T(tasks.ParseTaskFieldName(parsed.Get("sortBy")))
	}
	if parsed.Has("sortDesc") {
		/* Empty query parameter counts as "true" */
		if parsed.Get("sortDesc") == "" {
			taskQuery.SortDesc = true
		} else {
			newSortDesc, err := strconv.ParseBool(parsed.Get("sortDesc"))
			if err != nil {
				newSortDesc = false
			}
			taskQuery.SortDesc = newSortDesc
		}
	}
}

func (taskQuery TasksQuery[T]) String() string {
	var result []string
	/* Pagination */
	if taskQuery.Page != 1 {
		result = append(result, fmt.Sprintf("page=%d", taskQuery.Page))
	}
	if taskQuery.Size != defaultPageSize {
		result = append(result, fmt.Sprintf("size=%d", taskQuery.Size))
	}
	/* Search by */
	if taskQuery.SearchBy != "" {
		result = append(result, fmt.Sprintf("searchBy=%s", taskQuery.SearchBy))
	}
	/* Dates */
	defaultFromDate, defaultToDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	if taskQuery.FromDate != defaultFromDate.Format(utils.HTMLDateFormat) {
		result = append(result, fmt.Sprintf("from=%s", taskQuery.FromDate))
	}
	if taskQuery.ToDate != defaultToDate.Format(utils.HTMLDateFormat) {
		result = append(result, fmt.Sprintf("to=%s", taskQuery.ToDate))
	}
	/* Sorting */
	if taskQuery.SortBy != T(tasks.Datetime) {
		result = append(result, fmt.Sprintf("sortBy=%s", taskQuery.SortBy))
	}
	if taskQuery.SortDesc {
		result = append(result, "sortDesc")
	}

	if len(result) > 0 {
		return "?" + strings.Join(result, "&")
	} else {
		return ""
	}
}
