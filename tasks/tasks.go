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
	"todoer/jwt"
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
	task, err := Get(id)
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

func GetFromPayload(payload jwt.Payload) ([]Task, int, int) {
	result := slices.Clone(list)
	total := len(result)
	filter := payload.Filter
	page := payload.Page
	pageSize := payload.PageSize
	/* Filter */
	if filter != "" {
		result = slices.DeleteFunc(result, func(t Task) bool {
			return !strings.Contains(t.Description, filter)
		})
		total = len(result)
	}
	/* Sorting */
	switch payload.SortBy {
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
	if !payload.SortAsc {
		slices.Reverse(result)
	}

	/* Pagination */
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if page >= totalPages {
		page = totalPages - 1
	}
	/* Final result */
	startInd := pageSize * page
	endInd := min(startInd+pageSize, total)
	log.Printf("page %d: %d out of %d", page, pageSize, totalPages)
	return result[startInd:endInd], page, totalPages
}

func Get(id int) (Task, error) {
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
