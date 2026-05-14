package utils

import (
	"html/template"
	"time"
)

const HTMLDateFormat = "2006-01-02"

type SortableColumn int

const (
	Description SortableColumn = iota
	Date
)

/* Returns first and last day of current month */
func GetMonthBounds() (time.Time, time.Time) {
	now := time.Now()
	fromDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	toDate := fromDate.AddDate(0, 1, -1)
	return fromDate, toDate
}

/* FuncMap for HTML templates */
var TemplateFuncMap = template.FuncMap{
	/* Human-readable date formatting */
	"formatTime": func(t time.Time) string {
		return t.Format("2.01.2006 15:04:05")
	},
	/* Preset dates */
	"getFirstDayOfMonth": func() string {
		result, _ := GetMonthBounds()
		return result.Format(HTMLDateFormat)
	},
	"getLastDayOfMonth": func() string {
		_, result := GetMonthBounds()
		return result.Format(HTMLDateFormat)
	},
	"getFirstDayOfPreviousMonth": func() string {
		previousMonth := time.Now().AddDate(0, -1, 0)
		fromDate := time.Date(previousMonth.Year(), previousMonth.Month(), 1, 0, 0, 0, 0, previousMonth.Location())
		return fromDate.Format(HTMLDateFormat)
	},
	"getLastDayOfPreviousMonth": func() string {
		lastMonth := time.Now().AddDate(0, -1, 0)
		previousMonth := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())
		toDate := previousMonth.AddDate(0, 1, -1)
		return toDate.Format(HTMLDateFormat)
	},
	"getToday": func() string {
		now := time.Now()
		return now.Format(HTMLDateFormat)
	},
	"getYesterday": func() string {
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		return yesterday.Format(HTMLDateFormat)
	},
	"greet": func(name string) string {
		return "Hello, " + name + "!"
	},
}

func GetPagination(totalPages int, selectedPage int) []int {
	result := make([]int, 0, totalPages)
	for i := range totalPages {
		/* First 5 pages */
		if i < 5 {
			result = append(result, i+1)
			continue
		}
		/* Last 1 page */
		if i == totalPages-1 {
			result = append(result, i+1)
			continue
		}
		/* Around selectedPage */
		if i >= selectedPage-2 && i <= selectedPage {
			result = append(result, i+1)
			continue
		}
		/* Empty spaces */
		if result[len(result)-1] > 0 {
			result = append(result, 0)
		}
	}
	return result
}
