package tasks

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora"
	"libs.altipla.consulting/errors"
)

type taskStatus int

const (
	taskStatusPending = taskStatus(iota)
	taskStatusRunning
	taskStatusSuccess
	taskStatusFailed
)

type Task struct {
	Message string
	Handler func() error
	level   int

	mx     sync.RWMutex
	status taskStatus
	anim   animation
	err    error
}

func (task *Task) indentation() string {
	return strings.Repeat(" ", (task.level+1)*2)
}

func (task *Task) Run() error {
	task.mx.Lock()
	task.status = taskStatusRunning
	task.anim.restart()
	task.mx.Unlock()

	err := task.Handler()

	task.mx.Lock()
	defer task.mx.Unlock()
	task.status = taskStatusSuccess
	if err != nil {
		task.status = taskStatusFailed
		task.err = err
	}

	return errors.Trace(err)
}

func (task *Task) Render(w io.Writer) {
	task.mx.RLock()
	defer task.mx.RUnlock()

	fmt.Fprint(w, task.indentation())
	switch task.status {
	case taskStatusPending:
		fmt.Fprint(w, "  ")
	case taskStatusRunning:
		task.anim.render(w)
	case taskStatusFailed:
		fmt.Fprint(w, aurora.Red("✖ "))
	case taskStatusSuccess:
		fmt.Fprint(w, aurora.Green("✔ "))
	}

	fmt.Fprintln(w, task.Message)

	if task.status == taskStatusFailed {
		fmt.Fprintln(w, aurora.Yellow(task.indentation()+"  → ").Faint(), aurora.Red(task.err.Error()))
	}
}

var progressChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type animation struct {
	cur int
}

func (anim *animation) restart() {
	anim.cur = 0
}

func (anim *animation) render(w io.Writer) {
	fmt.Fprint(w, aurora.Yellow(progressChars[anim.cur]+" "))
	anim.cur = (anim.cur + 1) % len(progressChars)
}

type ParentTask struct {
	Message string
	Tasks   []*Task
	Disable bool

	mx     sync.RWMutex
	status taskStatus
}

func (task *ParentTask) Run() error {
	if task.Disable {
		return nil
	}

	task.mx.Lock()
	task.status = taskStatusRunning
	for _, child := range task.Tasks {
		child.level++
	}
	task.mx.Unlock()

	for _, child := range task.Tasks {
		if err := child.Run(); err != nil {
			task.mx.Lock()
			defer task.mx.Unlock()
			task.status = taskStatusFailed

			return errors.Trace(err)
		}
	}

	task.mx.Lock()
	defer task.mx.Unlock()
	task.status = taskStatusSuccess

	return nil
}

func (task *ParentTask) Render(w io.Writer) {
	if task.Disable {
		return
	}

	task.mx.RLock()
	defer task.mx.RUnlock()

	fmt.Fprint(w, "  ")
	switch task.status {
	case taskStatusPending:
		fmt.Fprint(w, "  ")
	case taskStatusRunning, taskStatusFailed:
		fmt.Fprint(w, aurora.Yellow("❯ "))
	case taskStatusSuccess:
		fmt.Fprint(w, aurora.Green("✔ "))
	}

	fmt.Fprintln(w, task.Message)

	switch task.status {
	case taskStatusRunning, taskStatusFailed:
		for _, child := range task.Tasks {
			child.Render(w)
		}
	}
}
