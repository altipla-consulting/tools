package tasks

import (
	"strings"
	"sync"
)

type logsCapturer struct {
	mx            sync.RWMutex
	lastLine      []byte
	ellipsisAdded bool
	startNewLine  bool
}

func (ct *logsCapturer) Write(buf []byte) (int, error) {
	ct.mx.Lock()
	defer ct.mx.Unlock()

	for _, b := range buf {
		// If a new character comes after a newline we start to fill up the newline.
		if ct.startNewLine && b != '\n' && b != '\r' {
			ct.startNewLine = false
			ct.ellipsisAdded = false
			ct.lastLine = nil
		}

		// Fill characters until a limit
		if len(ct.lastLine) < 60 {
			ct.lastLine = append(ct.lastLine, b)
		} else if !ct.ellipsisAdded {
			ct.ellipsisAdded = true
			ct.lastLine = append(ct.lastLine, '.', '.', '.')
		}

		if b == '\n' {
			ct.startNewLine = true
		}
	}

	return len(buf), nil
}

func (ct *logsCapturer) LastLine() string {
	ct.mx.RLock()
	defer ct.mx.RUnlock()

	if ct.lastLine == nil {
		return ""
	}
	return strings.TrimSpace(string(ct.lastLine))
}
