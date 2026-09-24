package utils

import (
	"net/http"
	"time"
	"todoer/server/token"
)

const (
	HTMLDateFormat            = "2006-01-02"
	NiceLookingDatetimeFormat = "2.01.2006 15:04:05"
)

type Payload struct {
	UserID string `json:"user_id"`
}

func GetTokenPayload(req *http.Request) Payload {
	return req.Context(). /* Get context from request */
				Value("token").(*token.Token[Payload]). /* Get "token" field */
				Payload
}

/* Returns first and last day of month */
func GetMonthBounds(year int, month time.Month) (time.Time, time.Time) {
	fromDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	toDate := fromDate.AddDate(0, 1, -1)
	return fromDate, toDate
}

func GetPagination(totalPages uint, selectedPage uint) []uint {
	result := make([]uint, 0, totalPages)
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
