package tasks

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	tasksFilename = "tasks.yaml"
	timeFormat = "2.01.2006 15:04:05"
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

func GetAll(filter string, size int, page int) []Task {
	/* Filter */
	result := slices.DeleteFunc(slices.Clone(list), func(t Task) bool {
		return !strings.Contains(t.Description, filter)
	})
	/* Sort */

	/* Pagination */
	totalPages := int(math.Ceil(float64(len(result)) / float64(size)))
	if page >= totalPages {
		page = totalPages - 1
	}
	/* Final result */
	startInd := size * page
	endInd := startInd + size
	if (endInd > len(result)) {
		endInd = len(result)
	}
	return result[startInd:endInd]
}

func Get(id int) (Task, error) {
	ind := slices.IndexFunc(list, func (t Task) bool { 
		return t.Id == id
	})
	if ind == -1 {
		return Task{}, fmt.Errorf("Task not found: %d", id)
	}
	return list[ind], nil
}

func Update(id int, newDescription string) error {
	ind := slices.IndexFunc(list, func (t Task) bool { 
		return t.Id == id
	})
	if ind == -1 {
		return fmt.Errorf("Task not found: %d", id)
	}
	list[ind].Description = newDescription
	return nil
}
