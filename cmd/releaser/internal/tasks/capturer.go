package tasks

import (
	"strings"
	"sync"
)

type Controller struct {
	mx            sync.RWMutex
	lastLine      []byte
	ellipsisAdded bool
	startNewLine  bool
	manualOutput  []string
}

func (ctrl *Controller) Write(buf []byte) (int, error) {
	ctrl.mx.Lock()
	defer ctrl.mx.Unlock()

	for _, b := range buf {
		// If a new character comes after a newline we start to fill up the newline.
		if ctrl.startNewLine && b != '\n' && b != '\r' {
			ctrl.startNewLine = false
			ctrl.ellipsisAdded = false
			ctrl.lastLine = nil
		}

		// Fill characters until a limit
		if len(ctrl.lastLine) < 60 {
			ctrl.lastLine = append(ctrl.lastLine, b)
		} else if !ctrl.ellipsisAdded {
			ctrl.ellipsisAdded = true
			ctrl.lastLine = append(ctrl.lastLine, '.', '.', '.')
		}

		if b == '\n' {
			ctrl.startNewLine = true
		}
	}

	return len(buf), nil
}

func (ctrl *Controller) LastLine() string {
	ctrl.mx.RLock()
	defer ctrl.mx.RUnlock()

	if ctrl.lastLine == nil {
		return ""
	}
	return strings.TrimSpace(string(ctrl.lastLine))
}

func (ctrl *Controller) AddManualOutput(lines ...string) {
	ctrl.mx.Lock()
	defer ctrl.mx.Unlock()
	ctrl.manualOutput = append(ctrl.manualOutput, lines...)
}

func (ctrl *Controller) ManualOutput() []string {
	ctrl.mx.RLock()
	defer ctrl.mx.RUnlock()
	return append([]string{}, ctrl.manualOutput...)
}
