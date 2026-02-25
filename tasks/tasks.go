package tasks

import (
	"github.com/goccy/go-yaml"
	"log"
	"os"
	"time"
)

const tasksFilename = "tasks.yaml"

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
	for i, t := range list {
		t, err := time.Parse(time.DateTime, t.Datetime)
		if err != nil {
			panic(err)
		}
		list[i].DatetimeParsed = t
		list[i].DatetimeFormatted = t.Format("2.01.2006 15:04:05")
	}
	log.Println("Tasks found:", len(list))
}

func Get() []Task {
	return list
}

func Save() {
}
