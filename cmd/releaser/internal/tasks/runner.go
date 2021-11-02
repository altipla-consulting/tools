package tasks

import (
	"io"
	"sync"
	"time"

	"github.com/gosuri/uilive"
	"libs.altipla.consulting/errors"
)

type Runner struct {
	Tasks []Runnable
}

type Runnable interface {
	Run() error
	Render(w io.Writer)
}

func (runner *Runner) Run() error {
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		out := uilive.New()
		for {
			for _, task := range runner.Tasks {
				task.Render(out)
			}
			out.Flush()
			time.Sleep(100 * time.Millisecond)

			select {
			case <-stop:
				for _, task := range runner.Tasks {
					task.Render(out)
				}
				out.Flush()
				return
			default:
			}
		}
	}()

	for _, task := range runner.Tasks {
		if err := task.Run(); err != nil {
			stop <- struct{}{}
			wg.Wait()
			return errors.Trace(err)
		}
	}

	stop <- struct{}{}
	wg.Wait()
	return nil
}
