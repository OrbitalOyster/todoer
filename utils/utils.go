package utils

import (
	"html/template"
	"strings"
	"time"
)

const (
	HTMLDateFormat            = "2006-01-02"
	NiceLookingDatetimeFormat = "2.01.2006 15:04:05"
)

type SortableField int

const (
	Description SortableField = iota
	Datetime
)

func (field SortableField) String() string {
	switch field {
	case Description:
		return "Description"
	case Datetime:
		return "Datetime"
	default:
		panic("Invalid SortableField")
	}
}

func ParseSortableField(field string) SortableField {
	switch strings.ToLower(field) {
	case "description":
		return Description
	case "datetime":
		return Datetime
	default:
		panic("Invalid SortableField")
	}
}

/* Returns first and last day of month */
func GetMonthBounds(year int, month time.Month) (time.Time, time.Time) {
	fromDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	toDate := fromDate.AddDate(0, 1, -1)
	return fromDate, toDate
}

/* FuncMap for HTML templates */
var TemplateFuncMap = template.FuncMap{
	/* Human-readable date formatting */
	"formatDatetime": func(t time.Time) string {
		return t.Format(NiceLookingDatetimeFormat)
	},
	/* Preset dates */
	"getFirstDayOfMonth": func() string {
		fromDate, _ := GetMonthBounds(time.Now().Year(), time.Now().Month())
		return fromDate.Format(HTMLDateFormat)
	},
	"getLastDayOfMonth": func() string {
		_, toDate := GetMonthBounds(time.Now().Year(), time.Now().Month())
		return toDate.Format(HTMLDateFormat)
	},
	"getFirstDayOfPreviousMonth": func() string {
		fromDate, _ := GetMonthBounds(time.Now().Year(), time.Now().AddDate(0, -1, 0).Month())
		return fromDate.Format(HTMLDateFormat)
	},
	"getLastDayOfPreviousMonth": func() string {
		_, toDate := GetMonthBounds(time.Now().Year(), time.Now().AddDate(0, -1, 0).Month())
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
	"parseSortableField": ParseSortableField,
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
