package tasks

import (
	"fmt"
	"strings"
)

type TaskStatus uint

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

func ParseStatus(status string) (TaskStatus, error) {
	switch strings.ToLower(status) {
	case "inprogress":
		return InProgress, nil
	case "done":
		return Done, nil
	case "failed":
		return Failed, nil
	default:
		return 0, fmt.Errorf("Invalid TaskStatus: %s", status)
	}
}
