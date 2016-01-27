package main

import (
	"fmt"
	"os"

	"github.com/casimir/gott"
)

func main() {
	tasks, err := gott.LoadFile(os.Getenv("HOME") + "/Dropbox/todo/todo.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d tasks\n", len(tasks))
}
