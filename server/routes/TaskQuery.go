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

type TasksQuery struct {
	Page     uint
	Size     uint
	SearchBy string
	FromDate time.Time
	ToDate   time.Time
	SortBy   tasks.TaskField
	SortDesc bool
}

func defaultTaskQuery() TasksQuery {
	fromDate, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	return TasksQuery{
		Page:     1,
		Size:     defaultPageSize,
		SearchBy: "",
		FromDate: fromDate,
		ToDate:   toDate,
		SortBy:   tasks.Datetime,
		SortDesc: false,
	}
}

func (taskQuery *TasksQuery) parse(rawQuery string) {
	parsed, err := url.ParseQuery(rawQuery)
	if err != nil {
		return
	}
	/* Page */
	if parsed.Has("page") {
		if parsedPage, err := strconv.Atoi(parsed.Get("page")); err == nil {
			taskQuery.Page = uint(parsedPage)
		}
	}
	/* Page size */
	if parsed.Has("size") {
		if parsedSize, err := strconv.Atoi(parsed.Get("size")); err == nil {
			taskQuery.Size = uint(parsedSize)
		}
	}
	/* Search by */
	if parsed.Has("searchBy") {
		taskQuery.SearchBy = parsed.Get("searchBy")
	}
	/* Dates */
	defaultFromDate, defaultToDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	if parsed.Has("from") {
		if fromDate, err := time.ParseInLocation(utils.HTMLDateFormat, parsed.Get("from"), time.Local); err != nil {
			taskQuery.FromDate = defaultFromDate
		} else {
			taskQuery.FromDate = fromDate
		}
	}
	if parsed.Has("to") {
		if toDate, err := time.ParseInLocation(utils.HTMLDateFormat, parsed.Get("to"), time.Local); err != nil {
			taskQuery.ToDate = defaultToDate
		} else {
			taskQuery.ToDate = toDate
		}
	}
	/* Sorting */
	if parsed.Has("sortBy") {
		sortByField, err := tasks.ParseTaskField(parsed.Get("sortBy"))
		if err != nil {
			panic("Invalid sortByField: " + sortByField.String())
		}
		taskQuery.SortBy = tasks.TaskField(sortByField)
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

func CreateQueryFromRequest(req *http.Request) (query TasksQuery, updated bool) {
	query = defaultTaskQuery()
	updated = false
	currentQuery := req.URL.RawQuery
	/* HTMX request, fish out the original URL */
	if hxCurrentUrl := req.Header.Get("HX-Current-URL"); hxCurrentUrl != "" {
		url, err := url.Parse(hxCurrentUrl)
		if err != nil {
			panic(err)
		}
		query.parse(url.RawQuery)
		currentQuery = query.String()
	}
	/* Parse the actual request */
	query.parse(req.URL.RawQuery)
	if currentQuery != query.String() {
		updated = true
	}
	/* "Naked" return */
	return
}

func (taskQuery TasksQuery) String() string {
	var fields []string
	/* Pagination */
	if taskQuery.Page != 1 {
		fields = append(fields, fmt.Sprintf("page=%d", taskQuery.Page))
	}
	if taskQuery.Size != defaultPageSize {
		fields = append(fields, fmt.Sprintf("size=%d", taskQuery.Size))
	}
	/* Search by */
	if taskQuery.SearchBy != "" {
		fields = append(fields, fmt.Sprintf("searchBy=%s", taskQuery.SearchBy))
	}
	/* Dates */
	defaultFromDate, defaultToDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	if taskQuery.FromDate != defaultFromDate {
		fields = append(fields, fmt.Sprintf("from=%s", taskQuery.FromDate.Format(utils.HTMLDateFormat)))
	}
	if taskQuery.ToDate != defaultToDate {
		fields = append(fields, fmt.Sprintf("to=%s", taskQuery.ToDate.Format(utils.HTMLDateFormat)))
	}
	/* Sorting */
	if taskQuery.SortBy != tasks.TaskField(tasks.Datetime) {
		fields = append(fields, fmt.Sprintf("sortBy=%s", taskQuery.SortBy))
	}
	if taskQuery.SortDesc {
		fields = append(fields, "sortDesc")
	}
	return strings.Join(fields, "&")
}
