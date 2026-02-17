package tasks

import (
	"github.com/goccy/go-yaml"
	"log"
	"os"
)

const tasksFile = "tasks.yaml"

type Task struct {
	Id          int    `yaml:"id"`
	User				string `yaml:"user"`
	Description string `yaml:"description"`
	Done        bool   `yaml:"done"`
}

var list []Task

func Load()  {
	log.Println("Loading tasks from", tasksFile)
	/* Load raw yaml */
	listRaw, err := os.ReadFile(tasksFile)
	if err != nil {
		panic(err)
	}
	/* Parse */
	if err := yaml.Unmarshal(listRaw, &list); err != nil {
		panic(err)
	}
	// list[0].Done = true
	// output, err := yaml.Marshal(tasks)
	// log.Println(string(output))
}

func Save() {
}
