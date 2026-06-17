package wad

import (
	"log"
	"os"
	"todoer/tasks"

	"github.com/goccy/go-yaml"
)

func Load(filename string) ([]User, []tasks.Task) {
	/* Load raw yaml */
	log.Println("Loading everything from", filename)
	rawData, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	/* Parse */
	var wad WAD
	if err := yaml.Unmarshal(rawData, &wad); err != nil {
		panic(err)
	}
	log.Println("Users found:", len(wad.Users))
	log.Println("Tasks found:", len(wad.Tasks))
	return wad.Users, wad.Tasks
}
