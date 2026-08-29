package tasks

import (
	"fmt"
	"strings"
)

type TaskStatus int

const (
	InProgress TaskStatus = iota
	Done
	Failed
)

func (status TaskStatus) String() string {
	switch status {
	case InProgress:
		return "InProgress"
	case Done:
		return "Done"
	case Failed:
		return "Failed"
	default:
		/* Major screwup */
		panic("Invalid TaskStatus")
	}
}

func ParseStatus(status string) TaskStatus {
	switch strings.ToLower(status) {
	case "inprogress":
		return InProgress
	case "done":
		return Done
	case "failed":
		return Failed
	default:
		panic(fmt.Sprintf("Invalid TaskStatus: %s", status))
	}
}
