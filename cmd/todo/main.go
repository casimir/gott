package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/casimir/gott"
	"github.com/casimir/gott/filters"
)

func main() {
	tasks, err := gott.LoadFile(os.Getenv("HOME") + "/Dropbox/todo/todo.txt")
	if err != nil {
		panic(err)
	}
	priotasks := tasks.Filter(filters.HasPrio)

	fmt.Printf("%d tasks, %d with priority\n", len(tasks), len(priotasks))
	projects := tasks.Projects()
	fmt.Printf("%d projects: %s\n", len(projects), strings.Join(projects, ", "))
	contexts := tasks.Contexts()
	fmt.Printf("%d contexts: %s\n", len(contexts), strings.Join(contexts, ", "))
	fmt.Println("----")
	for _, it := range priotasks {
		fmt.Println(it)
	}
}
