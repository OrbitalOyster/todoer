package utils

import (
	"time"
)

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
