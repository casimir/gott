package filters

import "github.com/casimir/gott"

// Done checks if the task is done.
func Done(td gott.TaskData) bool {
	return td.Done
}

// HasPrio checks if the task as a priority.
func HasPrio(td gott.TaskData) bool {
	return td.Priority != byte(0)
}

func FilterContext(context string) gott.FilterFn {
	return func(td gott.TaskData) bool {
		for _, it := range td.Context {
			if it == context {
				return true
			}
		}
		return false
	}
}

func FilterProject(project string) gott.FilterFn {
	return func(td gott.TaskData) bool {
		for _, it := range td.Project {
			if it == project {
				return true
			}
		}
		return false
	}
}
