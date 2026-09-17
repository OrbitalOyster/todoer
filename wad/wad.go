package wad

import (
	"log"
	"os"
	"strings"
	"todoer/tasks"
	"todoer/users"

	"github.com/goccy/go-yaml"
)

func Load(filename string) ([]users.User, []tasks.Task) {
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
	log.Println("Categories:", strings.Join(wad.Categories, ", "))
	log.Println("Tasks found:", len(wad.Tasks))
	return wad.Users, wad.Tasks
}
