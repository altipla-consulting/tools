package cli

import (
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
)

func Success(format string) {
	Successf(format)
}

func Successf(format string, args ...interface{}) {
	log.Println(emoji.Sprintf(":white_check_mark: "+format, args...))
}
