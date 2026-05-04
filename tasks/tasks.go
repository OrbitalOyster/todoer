package tasks

import (
	"cmp"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"todoer/utils"

	"github.com/goccy/go-yaml"
)

const (
	tasksFilename = "tasks.yaml"
	timeFormat    = "2.01.2006 15:04:05"
)

type Task struct {
	Id                int    `yaml:"id"`
	User              string `yaml:"user"`
	Datetime          string `yaml:"datetime"`
	DatetimeParsed    time.Time
	DatetimeFormatted string
	Description       string `yaml:"description"`
	Done              bool   `yaml:"done"`
}

var list []Task

func Check(idStr string) Task {
	id, err := strconv.Atoi(idStr)
	/* User sent stoopid */
	if err != nil {
		panic(err)
	}
	task, err := GetOne(id)
	/* No such task */
	if err != nil {
		panic(err)
	}
	return task
}

func Load() {
	log.Println("Loading tasks from", tasksFilename)
	/* Load raw yaml */
	listRaw, err := os.ReadFile(tasksFilename)
	if err != nil {
		panic(err)
	}
	/* Parse */
	if err := yaml.Unmarshal(listRaw, &list); err != nil {
		panic(err)
	}
	/* Format datetime */
	for i, task := range list {
		time, err := time.Parse(time.DateTime, task.Datetime)
		if err != nil {
			panic(err)
		}
		list[i].DatetimeParsed = time
		list[i].DatetimeFormatted = time.Format(timeFormat)
	}
	log.Println("Tasks found:", len(list))
}

func Get(fromDateStr string, toDateStr string, searchBy string, page int, pageSize int, sortBy utils.SortableColumn, sortAsc bool) ([]Task, int, int) {
	result := slices.Clone(list)
	/* Date */
	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	toDate, err := time.Parse("2006-01-02", toDateStr)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	result = slices.DeleteFunc(result, func(t Task) bool {
		/* "Not after 20/03/2026" means "Not after 20/03/2026 23:59:59"  */
		return t.DatetimeParsed.Before(fromDate) || t.DatetimeParsed.After(toDate.Add(time.Hour*24-time.Second))
	})
	/* Search */
	if searchBy != "" {
		result = slices.DeleteFunc(result, func(t Task) bool {
			return !strings.Contains(t.Description, searchBy)
		})
	}
	/* Number of tasks after all filtering */
	total := len(result)
	/* Nothing found - stop */
	if total == 0 {
		return nil, 1, 1
	}
	/* Sorting */
	switch sortBy {
	case utils.Description:
		slices.SortFunc(result, func(t1, t2 Task) int {
			return cmp.Compare(t1.Description, t2.Description)
		})
	case utils.Date:
		slices.SortFunc(result, func(t1, t2 Task) int {
			return t1.DatetimeParsed.Compare(t2.DatetimeParsed)
		})
	default:
	}
	/* On reverse order */
	if !sortAsc {
		slices.Reverse(result)
	}
	/* Pagination */
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if page >= totalPages {
		page = totalPages
	}
	if page <= 0 {
		page = 1
	}
	/* Final result */
	startInd := pageSize * (page - 1)
	endInd := min(startInd+pageSize, total)
	return result[startInd:endInd], totalPages, page
}

func GetOne(id int) (Task, error) {
	ind := slices.IndexFunc(list, func(t Task) bool {
		return t.Id == id
	})
	if ind == -1 {
		return Task{}, fmt.Errorf("Task not found: %d", id)
	}
	return list[ind], nil
}

func Update(id int, newDescription string) error {
	ind := slices.IndexFunc(list, func(t Task) bool {
		return t.Id == id
	})
	if ind == -1 {
		return fmt.Errorf("Task not found: %d", id)
	}
	list[ind].Description = newDescription
	return nil
}
