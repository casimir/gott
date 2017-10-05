package main

import (
	"fmt"
	"os"

	"github.com/casimir/gott"
	"github.com/casimir/gott/filters"
)

func main() {
	tasks, err := gott.LoadFile(os.Getenv("HOME") + "/Dropbox/todo/todo.txt")
	if err != nil {
		panic(err)
	}
	priotasks := tasks.Filter(filters.HasPrio)
	if len(priotasks) > 1 {
		fmt.Println(priotasks[0])
	} else {
		fmt.Println("No task.")
	}
}
