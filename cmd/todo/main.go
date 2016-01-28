package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/casimir/gott"
)

func main() {
	tasks, err := gott.LoadFile(os.Getenv("HOME") + "/Dropbox/todo/todo.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d tasks\n", len(tasks))
	projects := tasks.Projects()
	fmt.Printf("%d projects: %s\n", len(projects), strings.Join(projects, ", "))
	contexts := tasks.Contexts()
	fmt.Printf("%d contexts: %s\n", len(contexts), strings.Join(contexts, ", "))
}
