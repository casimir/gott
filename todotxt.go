package gott

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

type TaskData struct {
	Done     bool
	Priority byte
	Date     time.Time
	Project  []string
	Context  []string
	Text     string
}

func (td TaskData) String() string {
	var parts []string
	if td.Done {
		parts = append(parts, "x")
	}
	if td.Priority != byte(0) {
		s := fmt.Sprintf("(%c)", td.Priority)
		parts = append(parts, s)
	}
	if td.Date != (time.Time{}) {
		s := td.Date.Format(dateFormat)
		parts = append(parts, s)
	}
	parts = append(parts, td.Text)
	if len(td.Project) > 0 {
		s := "+" + strings.Join(td.Project, " +")
		parts = append(parts, s)
	}
	if len(td.Context) > 0 {
		s := "@" + strings.Join(td.Context, " @")
		parts = append(parts, s)
	}
	return strings.Join(parts, " ")
}

type TaskList []TaskData

func (l TaskList) Len() int      { return len(l) }
func (l TaskList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l TaskList) Less(i, j int) bool {
	a, b := l[i], l[j]
	if a.Priority == b.Priority {
		return a.Date == b.Date || a.Date.Before(b.Date)
	} else if a.Priority != byte(0) || b.Priority != byte(0) {
		return b.Priority == byte(0)
	}
	return a.Priority < b.Priority
}

func (ts TaskList) Projects() []string {
	buf := map[string]bool{}
	for _, task := range ts {
		for _, project := range task.Project {
			buf[project] = true
		}
	}

	var projects []string
	for k := range buf {
		projects = append(projects, k)
	}
	sort.Strings(projects)
	return projects
}

func (ts TaskList) Contexts() []string {
	buf := map[string]bool{}
	for _, task := range ts {
		for _, context := range task.Context {
			buf[context] = true
		}
	}

	var contexts []string
	for k := range buf {
		contexts = append(contexts, k)
	}
	sort.Strings(contexts)
	return contexts
}

type FilterFn func(TaskData) bool

func (ts TaskList) Filter(fns ...FilterFn) TaskList {
	var ret TaskList
	for _, it := range ts {
		for _, fn := range fns {
			if fn(it) {
				ret = append(ret, it)
				break
			}
		}
	}
	return ret
}

func (ts TaskList) Exclude(fns ...FilterFn) TaskList {
	var ret TaskList
	for _, it := range ts {
		for _, fn := range fns {
			if !fn(it) {
				ret = append(ret, it)
				break
			}
		}
	}
	return ret
}

type tokFn func(string, *TaskData) bool

type parserCtx struct {
	toks []string
	idx  int
}

func (ctx *parserCtx) get() (string, bool) {
	if ctx.idx >= len(ctx.toks) {
		return "", false
	}
	return ctx.toks[ctx.idx], true
}

func (ctx *parserCtx) parse(fn tokFn, td *TaskData) bool {
	tok, again := ctx.get()
	if !again || !fn(tok, td) {
		return false
	}
	ctx.idx++
	return true
}

func tokDone(tok string, td *TaskData) bool {
	if tok != "x" {
		return false
	}
	td.Done = true
	return true
}

func tokPriority(tok string, td *TaskData) bool {
	if len(tok) != 3 || tok[0] != '(' || tok[2] != ')' {
		return false
	}
	td.Priority = tok[1]
	return true
}

func tokDate(tok string, td *TaskData) bool {
	date, err := time.Parse(dateFormat, tok)
	if err != nil {
		return false
	}
	td.Date = date
	return true
}

func tokProject(tok string, td *TaskData) bool {
	if len(tok) < 2 || tok[0] != '+' {
		return false
	}
	td.Project = append(td.Project, tok[1:])
	return true
}

// Official spec says a priority can be `[A-Z]` but let's use `\w` instead.
func tokContext(tok string, td *TaskData) bool {
	if len(tok) < 2 || tok[0] != '@' {
		return false
	}
	td.Context = append(td.Context, tok[1:])
	return true
}

func tokText(tok string, td *TaskData) bool {
	if len(td.Text) > 0 {
		tok = " " + tok
	}
	td.Text += tok
	return true
}

func NewTask(input string) (td TaskData) {
	ctx := parserCtx{toks: strings.Fields(input)}
	ctx.parse(tokDone, &td)
	ctx.parse(tokPriority, &td)
	ctx.parse(tokDate, &td)
	for {
		if _, again := ctx.get(); !again {
			break
		}
		for _, fn := range []tokFn{tokProject, tokContext, tokText} {
			if ctx.parse(fn, &td) {
				break
			}
		}
	}
	return
}

func LoadFile(filename string) (TaskList, error) {
	var tasks TaskList
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return tasks, err
	}
	for _, line := range bytes.Split(raw, []byte("\n")) {
		tasks = append(tasks, NewTask(string(line)))
	}
	sort.Sort(tasks)
	return tasks, nil
}
