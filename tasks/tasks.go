package tasks

import (
	"github.com/goccy/go-yaml"
	"log"
	"os"
)

const tasksFilename = "tasks.yaml"

type Task struct {
	Id          int    `yaml:"id"`
	User				string `yaml:"user"`
	Description string `yaml:"description"`
	Done        bool   `yaml:"done"`
}

var list []Task

func Load()  {
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
	log.Println("Tasks found:", len(list))
}

func Save() {
}
