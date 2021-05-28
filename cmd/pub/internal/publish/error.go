package publish

import (
	"fmt"

	"libs.altipla.consulting/errors"
)

type PublishError struct {
	msg string
}

func NewPublishError(msg string) error {
	return errors.Trace(PublishError{msg})
}

func NewPublishErrorf(format string, args ...interface{}) error {
	return errors.Trace(PublishError{
		msg: fmt.Sprintf(format, args...),
	})
}

func (err PublishError) Error() string {
	return err.msg
}
